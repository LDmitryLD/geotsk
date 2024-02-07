package controller

import (
	"context"
	"encoding/json"
	"log"
	"projects/LDmitryLD/geotask/module/courierfacade/service"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type CourierController struct {
	courierService service.CourierFacer
}

func NewCourierController(courierService service.CourierFacer) *CourierController {
	return &CourierController{courierService: courierService}
}

func (c *CourierController) GetStatus(ctx *gin.Context) {
	// установить задержку в 50 миллисекунд
	time.Sleep(50 * time.Millisecond)

	// получить статус курьера из сервиса courierService используя метод GetStatus
	status := c.courierService.GetStatus(ctx)

	// отправить статус курьера в ответ
	ctx.JSON(200, status)

}

func (c *CourierController) MoveCourier(m webSocketMessage) {
	var cm CourierMove
	var err error

	// получить данные из m.Data и десериализовать их в структуру CourierMove

	data, ok := m.Data.([]byte)
	if !ok {
		log.Println("ошибка при приведении типа:", reflect.TypeOf(m.Data))
		return
	}

	if err = json.Unmarshal(data, &cm); err != nil {
		log.Println("ошибка при анмаршалинге MoveCourier():", err)
		return
	}

	ctx := context.Background()
	// вызвать метод MoveCourier у courierService
	c.courierService.MoveCourier(ctx, cm.Direction, cm.Zoom)
}
