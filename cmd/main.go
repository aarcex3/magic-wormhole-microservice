package main

import (
	"github.com/aarcex3/magic-wormhole-microservice/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", routes.Home)
	router.GET("/about", routes.About)
	router.POST("/send", routes.Send)
	router.GET("/recieve", routes.Receive)
	router.GET("/health", routes.Health)
	router.Run(":8000")
}
