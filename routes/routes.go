package routes

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aarcex3/magic-wormhole-microservice/models"
	"github.com/gin-gonic/gin"
	"github.com/psanford/wormhole-william/wormhole"
	qrcode "github.com/skip2/go-qrcode"
)

var client wormhole.Client

func monitorSendStatus(status <-chan wormhole.SendResult) {
	s := <-status
	if s.Error != nil {
		log.Printf("Send error: %s\n", s.Error)
	} else if s.OK {
		log.Println("Message sent successfully!")
	} else {
		log.Println("Send status not OK but no error reported.")
	}
}

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

	go monitorSendStatus(status)

	content := fmt.Sprintf("wormhole-transfer:%s", code)

	var png []byte
	png, err = qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

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
	if msg.Type != wormhole.TransferText {
		log.Fatalf("Expected a text message but got type %s", msg.Type)
	}
	msgBody, err := io.ReadAll(msg)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{"message": string(msgBody)})
}

func INDEX(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "Running",
	})

}
