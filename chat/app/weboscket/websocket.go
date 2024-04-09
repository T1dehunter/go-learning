package weboscket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocket struct {
	connections     map[*websocket.Conn]bool
	userAuthHandler func(message UserAuthMessage)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func NewWebSocket() *WebSocket {
	return &WebSocket{connections: make(map[*websocket.Conn]bool)}
}

func (websocket *WebSocket) Listen(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		websocket.HandleMessage(message)

		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}
}

func (websocket *WebSocket) HandleMessage(messageData []byte) {
	var message UserAuthMessage
	err := json.Unmarshal(messageData, &message)
	if err == nil && message.Name == "user_auth" {
		if websocket.userAuthHandler != nil {
			websocket.userAuthHandler(message)
		}
	}
	if err != nil {
		log.Println("Error unmarshalling message", err)
	}
	fmt.Println("Received message", message)
}

func (websocket *WebSocket) SubscribeOnUserAuth(handler func(message UserAuthMessage)) {
	if websocket.userAuthHandler == nil {
		websocket.userAuthHandler = handler
	}
}
