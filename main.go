package main

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/go-gcm"
	"github.com/labstack/echo"
)

func main() {

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.POST("/api/push-notif", SendGMToClient)
	e.Logger.Fatal(e.Start(":8080"))

}

// SendGMToClient is a function that will push a message to client
func SendGMToClient(c echo.Context) error {
	serverKey := "AAAAorSWiIM:APA91bGFfAnMlIt20vocPKeNkQc1qrblrUT6Q1AgAtY4ZyV4howzavhKrgtIBzFHi89i0b2Z62qcOy6xQsKpcNpl3MsX98UkkbbP51vNcz5LRtno5Dv737rOjXgUjxmjrvWGJk-5djVl"
	notification := gcm.Notification{
		Title:       c.FormValue("title"),
		Body:        c.FormValue("body"),
		ClickAction: c.FormValue("clickAction"),
	}
	msg := gcm.HttpMessage{
		Data:            map[string]interface{}{"message": c.FormValue("message")},
		RegistrationIds: []string{c.FormValue("client_token")},
		Notification:    &notification,
	}
	response, err := gcm.SendHttp(serverKey, msg)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Response ", response)
		fmt.Println("Response ", response.Success)
		fmt.Println("MessageID ", response.MessageId)
		fmt.Println("Failure ", response.Failure)
		fmt.Println("Error ", response.Error)
		fmt.Println("Results ", response.Results)
	}

	t := time.Now()
	uuid, errUUID := newUUID()
	if errUUID != nil {
		fmt.Printf("error: %v\n", err)
	}

	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"requestID":   uuid,
			"now":         t.Format("2006/01/02 15:04:05"),
			"code":        strconv.Itoa(http.StatusBadRequest) + "02",
			"err message": err.Error(),
			"data":        "[]",
		})
	} else {
		return c.JSON(http.StatusOK, echo.Map{
			"requestID":   uuid,
			"now":         t.Format("2006/01/02 15:04:05"),
			"code":        strconv.Itoa(http.StatusOK) + "01",
			"err message": response.Error,
			"data":        response.Results,
		})

	}
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
