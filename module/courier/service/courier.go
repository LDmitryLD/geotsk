package service

import (
	"context"
	"math"
	"projects/LDmitryLD/geotask/geo"
	"projects/LDmitryLD/geotask/module/courier/models"
	"projects/LDmitryLD/geotask/module/courier/storage"

	"github.com/redis/go-redis/v9"
)

const (
	DirectionUp    = 0
	DirectionDown  = 1
	DirectionLeft  = 2
	DirectionRight = 3
)

const (
	DefaultCourierLat = 59.9311
	DefaultCourierLng = 30.3609
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Courierer
type Courierer interface {
	GetCourier(ctx context.Context) (*models.Courier, error)
	MoveCourier(courier models.Courier, direction, zoom int) (models.Courier, error)
	Save(ctx context.Context, courier models.Courier) error
}

type CourierService struct {
	courierStorage storage.CourierStorager
	allowedZone    geo.PolygonChecker
	disabledZone   []geo.PolygonChecker
}

func NewCourierService(courierStorage storage.CourierStorager, allowedZone geo.PolygonChecker, disabledZone []geo.PolygonChecker) Courierer {
	return &CourierService{courierStorage: courierStorage, allowedZone: allowedZone, disabledZone: disabledZone}
}

func (c *CourierService) GetCourier(ctx context.Context) (*models.Courier, error) {
	// получаем курьера из хранилища используя метод GetOne из storage/courier.go
	var courier *models.Courier
	var err error

	courier, err = c.courierStorage.GetOne(ctx)
	if err == redis.Nil {
		courier = &models.Courier{
			Location: models.Point{
				Lat: DefaultCourierLat,
				Lng: DefaultCourierLng,
			},
		}
	} else if err != nil {
		return nil, err
	}

	// проверяем, что курьер находится в разрешенной зоне
	// если нет, то перемещаем его в случайную точку в разрешенной зоне

	courierPoint := geo.Point{
		Lat: courier.Location.Lat,
		Lng: courier.Location.Lng,
	}

	if !c.allowedZone.Contains(courierPoint) {
		randPoint := c.allowedZone.RandomPoint()
		courier.Location.Lat = randPoint.Lat
		courier.Location.Lng = randPoint.Lng
	}

	// сохраняем новые координаты курьера
	if err := c.courierStorage.Save(ctx, *courier); err != nil {
		return nil, err
	}

	return courier, nil
}

// MoveCourier : direction - направление движения курьера, zoom - зум карты
func (c *CourierService) MoveCourier(courier models.Courier, direction, zoom int) (models.Courier, error) {
	// точность перемещения зависит от зума карты использовать формулу 0.001 / 2^(zoom - 14)
	// 14 - это максимальный зум карты
	precision := 0.001 / math.Pow(2, float64(zoom-14))

	switch direction {
	case DirectionUp:
		courier.Location.Lat += precision
	case DirectionDown:
		courier.Location.Lat -= precision
	case DirectionLeft:
		courier.Location.Lng -= precision
	case DirectionRight:
		courier.Location.Lng += precision
	}

	// далее нужно проверить, что курьер не вышел за границы зоны
	// если вышел, то нужно переместить его в случайную точку внутри зоны
	if !c.allowedZone.Contains(geo.Point(courier.Location)) {
		randPoint := c.allowedZone.RandomPoint()
		courier.Location.Lat = randPoint.Lat
		courier.Location.Lng = randPoint.Lng
	}

	// далее сохранить изменения в хранилище
	err := c.courierStorage.Save(context.Background(), courier)
	return courier, err
}

func (c *CourierService) Save(ctx context.Context, courier models.Courier) error {
	return c.courierStorage.Save(ctx, courier)
}
