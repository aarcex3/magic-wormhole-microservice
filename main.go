package main

import (
	"github.com/aarcex3/magic-wormhole-microservice/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", routes.INDEX)

	router.POST("/send", routes.SEND)

	router.GET("/recieve", routes.RECIEVE)

	router.Run()
}
