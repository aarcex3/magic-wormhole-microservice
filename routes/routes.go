package routes

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aarcex3/magic-wormhole-microservice/models"
	"github.com/aarcex3/magic-wormhole-microservice/utils"
	"github.com/aarcex3/magic-wormhole-microservice/views"
	"github.com/gin-gonic/gin"
	"github.com/psanford/wormhole-william/wormhole"
)

var client wormhole.Client

func Send(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		handleError(c, http.StatusBadRequest, "Bad request", err)
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer file.Close()

	extension := utils.GetExtension(fileHeader)
	tempFile, err := utils.CreateTempFile(file, extension)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer os.Remove(tempFile.Name())

	code, status, err := client.SendFile(c, filepath.Base(tempFile.Name()), tempFile)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	go utils.MonitorStatus(status)

	url := utils.GenerateURL(code)
	png := utils.GenerateQR(url)

	c.Header("Content-Type", "image/png")
	c.Data(http.StatusOK, "image/png", png)
}

func Receive(c *gin.Context) {
	var code models.Code

	if err := c.BindJSON(&code); err != nil {
		handleError(c, http.StatusBadRequest, "Bad request", err)
		return
	}

	file, err := client.Receive(c, code.Code)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	tempFile, err := os.CreateTemp("./temp", "empfile_")
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}
	defer os.Remove(tempFile.Name())

	_, err = io.Copy(tempFile, file)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(tempFile.Name()))
	c.Header("Content-Type", "application/octet-stream")
	c.File(tempFile.Name())
}

func Home(c *gin.Context) {
	c.Status(200)
	utils.Render(c, views.Index())
}

func About(c *gin.Context) {
	c.Status(200)
	utils.Render(c, views.About())
}

func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Running"})
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, gin.H{"error": message})
	log.Println("Error:", err)
}
