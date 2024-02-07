package service

import (
	"context"
	"log"
	cservice "projects/LDmitryLD/geotask/module/courier/service"
	cfm "projects/LDmitryLD/geotask/module/courierfacade/models"
	"projects/LDmitryLD/geotask/module/order/models"
	oservice "projects/LDmitryLD/geotask/module/order/service"
)

const (
	CourierVisibilityRadius = 2800
	CourierGetOrderRadius   = 5
	unitMeters              = "m"
)

type CourierFacer interface {
	MoveCourier(ctx context.Context, direction, zoom int)
	GetStatus(ctx context.Context) cfm.CourierStatus
}

type CourierFacade struct {
	courierService cservice.Courierer
	orderService   oservice.Orderer
}

func NewCourierFacade(courierService cservice.Courierer, orderService oservice.Orderer) CourierFacer {
	return &CourierFacade{courierService: courierService, orderService: orderService}
}

func (c *CourierFacade) MoveCourier(ctx context.Context, direction int, zoom int) {
	courier, err := c.courierService.GetCourier(ctx)
	if err != nil {
		return
	}

	movedCourier, err := c.courierService.MoveCourier(*courier, direction, zoom)
	if err != nil {
		log.Println("ошибка при передвижении курьера:", err)
	}

	lat := movedCourier.Location.Lat
	lng := movedCourier.Location.Lng

	orders, err := c.orderService.GetByRadius(ctx, lng, lat, CourierGetOrderRadius, unitMeters)
	if err != nil {
		log.Println("ошибка при получении заказов в радиусе.")
	}
	if len(orders) == 0 {
		return
	}

	var order models.Order
	for _, ord := range orders {
		if !order.IsDelivered {
			order = ord
			break
		}
	}

	movedCourier.Score++

	if err = c.orderService.RemoveOrder(ctx, order); err != nil {
		log.Println("ошибка при удалении заказа:", err)
	}

	if err = c.courierService.Save(ctx, movedCourier); err != nil {
		log.Println("ошибка при сохранении курьера:", err)
	}

}

func (c *CourierFacade) GetStatus(ctx context.Context) cfm.CourierStatus {
	courier, err := c.courierService.GetCourier(ctx)
	if err != nil {
		log.Println("ошибка при получении курьера из кэша:", err)
		return cfm.CourierStatus{}
	}
	lat := courier.Location.Lat
	lng := courier.Location.Lng

	orders, err := c.orderService.GetByRadius(ctx, lng, lat, CourierVisibilityRadius, unitMeters)
	if err != nil {
		log.Println("ошибка при получении заказов в радиусе:", err)
	}

	status := cfm.CourierStatus{
		Courier: *courier,
		Orders:  orders,
	}

	return status
}
