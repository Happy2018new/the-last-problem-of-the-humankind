package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/", CheckAlive)
	router.GET("/process_exit", ProcessExist)
	router.POST("/place_nbt_block", PlaceNBTBlock)
	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})

	return router
}

func RunServer() {
	router := InitRouter()
	router.Run(fmt.Sprintf(":%d", *standardServerPort))
}
