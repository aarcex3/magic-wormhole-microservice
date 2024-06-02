package utils

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strings"

	"github.com/psanford/wormhole-william/wormhole"
	qrcode "github.com/skip2/go-qrcode"
)

func MonitorStatus(status <-chan wormhole.SendResult) {
	s := <-status
	if s.Error != nil {
		log.Printf("Send error: %s\n", s.Error)
	} else if s.OK {
		log.Println("Message sent successfully!")
	} else {
		log.Println("Send status not OK but no error reported.")
	}
}

func GenerateQR(content string) []byte {
	var png []byte
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}
	return png
}

func GenerateURL(code string) string {
	return fmt.Sprintf("wormhole-transfer:%s", code)
}

func GetExtension(file *multipart.FileHeader) string {
	return strings.Split(file.Filename, ".")[len(strings.Split(file.Filename, "."))-1]
}

func CreateTempFile(file io.Reader, extension string) (*os.File, error) {
	pattern := fmt.Sprintf("upload-*.%s", extension)
	return os.CreateTemp("./temp", pattern)
}
