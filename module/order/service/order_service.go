package service

import (
	"context"
	"math"
	"math/rand"
	"projects/LDmitryLD/geotask/geo"
	"projects/LDmitryLD/geotask/module/order/models"
	"projects/LDmitryLD/geotask/module/order/storage"
	"time"
)

const (
	minDeliveryPrice = 100.00
	maxDeliveryPrice = 500.00

	maxOrderPrice = 3000.00
	minOrderPrice = 1000.00

	orderMaxAge = 2 * time.Minute
)

//go:generate go run github.com/vektra/mockery/v2@v2.35.4 --name=Orderer
type Orderer interface {
	GetByRadius(ctx context.Context, lng, lat, radius float64, unit string) ([]models.Order, error)
	Save(ctx context.Context, order models.Order) error
	GetCount(ctx context.Context) (int, error)
	RemoveOldOrders(ctx context.Context) error
	GenerateOrder(ctx context.Context) error
	RemoveOrder(ctx context.Context, order models.Order) error
}

type OrderService struct {
	storage       storage.OrderStorager
	allowedZone   geo.PolygonChecker
	disabledZones []geo.PolygonChecker
}

func NewOrderService(storage storage.OrderStorager, allowedZone geo.PolygonChecker, disallowedZone []geo.PolygonChecker) Orderer {
	return &OrderService{storage: storage, allowedZone: allowedZone, disabledZones: disallowedZone}
}

func (o *OrderService) GetByRadius(ctx context.Context, lng float64, lat float64, radius float64, unit string) ([]models.Order, error) {
	return o.storage.GetByRadius(ctx, lng, lat, radius, unit)
}

func (o *OrderService) Save(ctx context.Context, order models.Order) error {
	return o.storage.Save(ctx, order, orderMaxAge)
}

func (o *OrderService) GetCount(ctx context.Context) (int, error) {
	return o.storage.GetCount(ctx)
}

func (o *OrderService) RemoveOldOrders(ctx context.Context) error {
	return o.storage.RemoveOldOrders(ctx, orderMaxAge)
}

func (o *OrderService) GenerateOrder(ctx context.Context) error {
	id, err := o.storage.GenerateUniqueID(ctx)
	if err != nil {
		return err
	}
	rand.Seed(time.Now().UnixNano())

	price := minOrderPrice + rand.Float64()*(maxOrderPrice-minOrderPrice)
	deliveryPrice := minDeliveryPrice + rand.Float64()*(maxDeliveryPrice-minDeliveryPrice)

	randPoint := geo.GetRandomAllowedLocation(o.allowedZone, o.disabledZones)

	order := models.Order{
		ID:            id,
		Price:         math.Round(price*100) / 100,
		DeliveryPrice: math.Round(deliveryPrice*100) / 100,
		Lat:           randPoint.Lat,
		Lng:           randPoint.Lng,
		IsDelivered:   false,
		CreatedAt:     time.Now(),
	}

	return o.Save(ctx, order)
}

func (o *OrderService) RemoveOrder(ctx context.Context, order models.Order) error {
	return o.storage.RemoveOrder(ctx, order)
}
