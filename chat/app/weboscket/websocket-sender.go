package weboscket

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsSender struct {
	connection *websocket.Conn
}

func NewWsSender(conn *websocket.Conn) *WsSender {
	return &WsSender{connection: conn}
}

func (wsSender *WsSender) SendMessageToUser(connectionID int, message string) {
	if err := wsSender.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("Error sending message", err)
		return
	}
}

func (wsSender *WsSender) SendMessageToRoom(connectionID int, message string) {
	if err := wsSender.connection.Close(); err != nil {
		log.Println("Error closing connection", err)
		return
	}
}
