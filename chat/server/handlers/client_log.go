package handlers

import (
	"chat/server/weboscket"
	"fmt"
	"time"
)

func HandleClientLogMsg(message weboscket.ClientLogMsg, response *weboscket.Response) {
	fmt.Printf("Handler HandleClientLogMsg received message -> %+v\n", message)
	response.SendLog(weboscket.LogMsg{
		Type:      "client",
		Title:     message.Payload.Text,
		CreatedAt: time.Now().Format("15:04:05"),
	})
}
