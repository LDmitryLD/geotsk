package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"projects/LDmitryLD/geotask/module/order/models"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=OrderStorager
type OrderStorager interface {
	Save(ctx context.Context, order models.Order, maxAge time.Duration) error
	GetByID(ctx context.Context, orderID int) (*models.Order, error)
	GenerateUniqueID(ctx context.Context) (int64, error)
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error)
	GetCount(ctx context.Context) (int, error)
	RemoveOldOrders(ctx context.Context, maxAge time.Duration) error
	RemoveOrder(ctx context.Context, order models.Order) error
}

type OrderStorage struct {
	storage *redis.Client
}

func NewOrderStorage(storage *redis.Client) OrderStorager {
	return &OrderStorage{storage: storage}
}

func (o *OrderStorage) Save(ctx context.Context, order models.Order, maxAge time.Duration) error {
	return o.saveOrderWithGeo(ctx, order, maxAge)
}

func (o *OrderStorage) RemoveOldOrders(ctx context.Context, maxAge time.Duration) error {
	// получить ID всех старых ордеров, которые нужно удалить
	// используя метод ZRangeByScore
	// старые ордеры это те, которые были созданы две минуты назад
	// и более

	max := time.Now().Add(-maxAge).Unix()
	opt := &redis.ZRangeBy{
		Max: fmt.Sprintf("%d", max),
		Min: "0",
	}

	oldOrders, err := o.storage.ZRangeByScore(ctx, "orders", opt).Result()
	if err != nil {
		log.Println("ошибка при получении старых заказов:", err)
		return err
	}

	// Проверить количество старых ордеров
	// удалить старые ордеры из redis используя метод ZRemRangeByScore где ключ "orders" min "-inf" max "(время создания старого ордера)"
	// удалять ордера по ключу не нужно, они будут удалены автоматически по истечению времени жизни

	if len(oldOrders) > 0 {
		key := "orders"
		c, err := o.storage.ZRemRangeByScore(ctx, key, "-inf", opt.Max).Result()
		if err != nil {
			log.Println("ошибка при удалении заказов:", err)
			return err
		}
		log.Println("заказ удалён ", c)
	}

	return nil
}

func (o *OrderStorage) RemoveOrder(ctx context.Context, order models.Order) error {
	key := "orders"
	time := fmt.Sprintf("%d", order.CreatedAt.Unix())
	_, err := o.storage.ZRemRangeByScore(ctx, key, time, time).Result()
	if err != nil {
		return err
	}

	key = fmt.Sprintf("order:%d", order.ID)

	_, err = o.storage.Del(ctx, key).Result()
	if err != nil {
		return err
	}

	return nil
}

func (o *OrderStorage) GetByID(ctx context.Context, orderID int) (*models.Order, error) {
	var err error
	var data []byte
	var order models.Order

	// получаем ордер из redis по ключу order:ID
	key := fmt.Sprintf("order:%d", orderID)
	res, err := o.storage.Get(ctx, key).Result()
	// проверяем что ордер не найден исключение redis.Nil, в этом случае возвращаем nil, nil
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		log.Println("ошибка при получении заказа по ID:", err)
		return nil, err
	}

	data = []byte(res)
	// десериализуем ордер из json
	if err = json.Unmarshal(data, &order); err != nil {
		log.Println("ошибка при анмаршалинге заказа по ID:", err)
		return nil, err
	}

	return &order, nil
}

func (o *OrderStorage) saveOrderWithGeo(ctx context.Context, order models.Order, maxAge time.Duration) error {
	var err error
	var data []byte

	// сериализуем ордер в json
	data, err = json.Marshal(order)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("order:%d", order.ID)

	// сохраняем ордер в json redis по ключу order:ID с временем жизни maxAge
	if err := o.storage.Set(ctx, key, data, maxAge).Err(); err != nil {
		log.Println("ошибка при сохранении заказа в Redis: ", err)
		return err
	}

	// добавляем ордер в гео индекс используя метод GeoAdd где Name - это ключ ордера, а Longitude и Latitude - координаты
	err = o.storage.GeoAdd(ctx, "orders_location", &redis.GeoLocation{
		Name:      key,
		Latitude:  order.Lat,
		Longitude: order.Lng,
	}).Err()
	if err != nil {
		log.Println("ошибка при добавлении заказа в гео индекс: ", err)
		return err
	}

	// zset сохраняем ордер для получения количества заказов со сложностью O(1)
	// Score - время создания ордера
	err = o.storage.ZAdd(ctx, "orders", redis.Z{Score: float64(order.CreatedAt.Unix()), Member: key}).Err()
	if err != nil {
		log.Println("ошибка при ZAdd: ", err)
		return err
	}

	return nil
}

func (o *OrderStorage) GetCount(ctx context.Context) (int, error) {
	// получить количество ордеров в упорядоченном множестве используя метод ZCard
	count, err := o.storage.ZCard(ctx, "orders").Result()
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (o *OrderStorage) GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error) {
	var err error
	var orders []models.Order
	//var data []byte
	var ordersLocation []redis.GeoLocation

	// используем метод getOrdersByRadius для получения ID заказов в радиусе
	ordersLocation, err = o.getOrdersByRadius(ctx, lng, lat, radius, unit)
	// обратите внимание, что в случае отсутствия заказов в радиусе
	// метод getOrdersByRadius должен вернуть nil, nil (при ошибке redis.Nil)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	orders = make([]models.Order, 0, len(ordersLocation))
	// проходим по списку ID заказов и получаем данные о заказе
	// получаем данные о заказе по ID из redis по ключу order:ID

	for _, orderLoc := range ordersLocation {
		key := orderLoc.Name
		orderRaw, err := o.storage.Get(ctx, key).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			log.Println("ошибка при получении заказа из кэша: ", err)
			continue
		}

		var order models.Order
		if err := json.Unmarshal([]byte(orderRaw), &order); err != nil {
			continue
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (o *OrderStorage) getOrdersByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]redis.GeoLocation, error) {
	// в данном методе мы получаем список ордеров в радиусе от точки
	// возвращаем список ордеров с координатами и расстоянием до точки

	query := &redis.GeoRadiusQuery{
		Radius:      radius,
		Unit:        unit,
		WithCoord:   true,
		WithDist:    true,
		WithGeoHash: true,
	}

	return o.storage.GeoRadius(ctx, "orders_location", lng, lat, query).Result()
}

func (o *OrderStorage) GenerateUniqueID(ctx context.Context) (int64, error) {
	var err error
	var id int64

	// генерируем уникальный ID для ордера
	// используем для этого redis incr по ключу order:id

	id, err = o.storage.Incr(ctx, "order:id").Result()
	if err != nil {
		log.Println("ошибка при получении уникального ID: ", err)
		return 0, err
	}

	return id, nil
}
