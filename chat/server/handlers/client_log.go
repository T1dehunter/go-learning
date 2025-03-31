package handlers

import (
	"chat/server/weboscket"
	"time"
)

func HandleClientLogMsg(message weboscket.ClientLogMsg, response *weboscket.Response) {
	response.SendLog(weboscket.LogMsg{
		Type:      "client",
		Title:     message.Payload.Text,
		CreatedAt: time.Now().Format("15:04:05"),
	})
}
