package order

import (
	"context"
	"log"
	"projects/LDmitryLD/geotask/module/order/service"
	"time"
)

const (
	orderCleanInterval = 5 * time.Second
)

type OrderCleaner struct {
	orderService service.Orderer
}

func NewOrderCleaner(orderService service.Orderer) *OrderCleaner {
	return &OrderCleaner{orderService: orderService}
}

func (o *OrderCleaner) Run() {
	ticker := time.NewTicker(orderCleanInterval)
	ctx := context.Background()
	go func() {
		for {
			select {
			case <-ticker.C:
				err := o.orderService.RemoveOldOrders(ctx)
				if err != nil {
					log.Println("ошибка при удалении старых заказов:", err)
				}
			}
		}
	}()
}
