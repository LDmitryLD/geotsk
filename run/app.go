package run

import (
	"context"
	"log"
	"net/http"
	"os"
	"projects/LDmitryLD/geotask/cache"
	"projects/LDmitryLD/geotask/geo"
	cservice "projects/LDmitryLD/geotask/module/courier/service"
	cstorage "projects/LDmitryLD/geotask/module/courier/storage"
	"projects/LDmitryLD/geotask/module/courierfacade/controller"
	cfservice "projects/LDmitryLD/geotask/module/courierfacade/service"
	oservice "projects/LDmitryLD/geotask/module/order/service"
	ostorage "projects/LDmitryLD/geotask/module/order/storage"
	"projects/LDmitryLD/geotask/router"
	"projects/LDmitryLD/geotask/server"
	"projects/LDmitryLD/geotask/workers/order"
	"time"

	"github.com/gin-gonic/gin"
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	rclient := cache.NewRedisClient(host, port)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	_, err := rclient.Ping(ctx).Result()
	if err != nil {
		return err
	}

	allowedZone := geo.NewAllowedZone()

	disAllowedZones := []geo.PolygonChecker{geo.NewDisAllowedZone1(), geo.NewDisAllowedZone2()}

	orderStorage := ostorage.NewOrderStorage(rclient)
	orderService := oservice.NewOrderService(orderStorage, allowedZone, disAllowedZones)

	orderGenerator := order.NewOrderGenerator(orderService)
	orderGenerator.Run()

	oldOrderCleaner := order.NewOrderCleaner(orderService)
	oldOrderCleaner.Run()

	courierStorage := cstorage.NewCourierStorage(rclient)
	courierService := cservice.NewCourierService(courierStorage, allowedZone, disAllowedZones)

	courierFacade := cfservice.NewCourierFacade(courierService, orderService)

	courierController := controller.NewCourierController(courierFacade)

	routes := router.NewRouter(courierController)

	r := server.NewHTTPServer()

	api := r.Group("/api")

	routes.CourierAPI(api)

	mainRoute := r.Group("/")

	routes.Swagger(mainRoute)

	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("public"))))

	log.Println("SERVER STARTED")
	return r.Run()
}
