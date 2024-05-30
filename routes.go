package main

import (
	"io"
	"log"
	"net/http"

	"github.com/aarcex3/magic-wormhole-microservice/models"
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

	go func() {
		s := <-status
		if s.Error != nil {
			log.Printf("Send error: %s\n", s.Error)
		} else if s.OK {
			log.Println("OK!")
		} else {
			log.Println("Hmm not ok but also not error")
		}
	}()

	c.JSON(200, gin.H{"code": code})
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
