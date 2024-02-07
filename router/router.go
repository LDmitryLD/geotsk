package router

import (
	"projects/LDmitryLD/geotask/module/courierfacade/controller"

	"github.com/gin-gonic/gin"
)

type Router struct {
	courier *controller.CourierController
}

func NewRouter(courier *controller.CourierController) *Router {
	return &Router{courier: courier}
}

func (r *Router) CourierAPI(router *gin.RouterGroup) {
	router.GET("/status", r.courier.GetStatus)
	router.GET("/ws", r.courier.Websocket)
}

func (r *Router) Swagger(router *gin.RouterGroup) {
	router.GET("/swagger", swaggerUI)
	router.Static("/public/", "./public")
}
