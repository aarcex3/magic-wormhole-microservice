package routes

import (
	"io"
	"log"
	"net/http"

	"github.com/aarcex3/magic-wormhole-microservice/models"
	"github.com/aarcex3/magic-wormhole-microservice/utils"
	"github.com/gin-gonic/gin"
	"github.com/psanford/wormhole-william/wormhole"
)

var client wormhole.Client

func SEND(c *gin.Context) {
	var message models.Message

	if err := c.BindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		log.Println("Error:", err)
		return
	}

	code, status, err := client.SendText(c, message.Text)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		log.Println("Error sending message:", err)
		return
	}

	go utils.MonitorStatus(status)

	content := utils.GenerateURL(code)
	png := utils.GenerateQR(content)

	c.Header("Content-Type", "image/jpeg")
	c.Data(http.StatusOK, "image/png", png)
}

func RECIEVE(c *gin.Context) {
	var code models.Code

	if err := c.BindJSON(&code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
		log.Println("Error:", err)
		return
	}

	msg, err := client.Receive(c, code.Code)
	if err != nil {
		log.Fatal(err)
	}

	msgBody, err := io.ReadAll(msg)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": string(msgBody)})
}

func HOME(c *gin.Context) {

	c.HTML(http.StatusOK, "index.html", nil)

}

func HEALTH(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{"message": "Running"})

}
