package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ypli0629/yoyo/middlewares"
)

func main() {
	r := gin.New()
	// middlewares
	// logger
	r.Use(middlewares.Logger())
	// recovery
	r.Use(gin.Recovery())

	r.Run()
}
