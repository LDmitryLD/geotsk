package server

import "github.com/gin-gonic/gin"

func NewHTTPServer(opts ...func(group *gin.RouterGroup)) *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")

	for _, opt := range opts {
		opt(api)
	}

	return r
}
