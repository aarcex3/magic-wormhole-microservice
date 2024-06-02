package main

import (
	"github.com/aarcex3/magic-wormhole-microservice/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", routes.HOME)
	router.POST("/send", routes.SEND)
	router.GET("/recieve", routes.RECIEVE)
	router.GET("/health", routes.HEALTH)
	router.Run()
}
