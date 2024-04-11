package weboscket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketServer struct {
	connections           map[*websocket.Conn]bool
	userAuthHandler       func(message UserAuthMessage, wsSender WebsocketSender)
	userJoinToRoomHandler func(message UserJoinToRoomMessage, wsSender WebsocketSender)
}

func NewWebSocketServer() *WebSocketServer {
	return &WebSocketServer{connections: make(map[*websocket.Conn]bool)}
}

func (wsServer *WebSocketServer) Listen(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		wsServer.HandleMessage(conn, message)
	}
}

func (wsServer *WebSocketServer) HandleMessage(conn *websocket.Conn, messageData []byte) {
	var message UserAuthMessage
	err := json.Unmarshal(messageData, &message)
	if err == nil && message.Name == "user_auth" {
		if wsServer.userAuthHandler != nil {
			sender := newWsSender(conn)
			wsServer.userAuthHandler(message, sender)
		}
	}

	var userJoinToRoomMessage UserJoinToRoomMessage
	err = json.Unmarshal(messageData, &userJoinToRoomMessage)
	if err == nil && userJoinToRoomMessage.Name == "user_join_to_room" {
		if wsServer.userJoinToRoomHandler != nil {
			sender := NewWsSender(conn)
			wsServer.userJoinToRoomHandler(userJoinToRoomMessage, sender)
		}
	}

	if err != nil {
		log.Println("Error unmarshalling message", err)
	}
	fmt.Println("Received message", message)
}

func (wsServer *WebSocketServer) SubscribeOnUserAuth(handler func(message UserAuthMessage, ws WebsocketSender)) {
	if wsServer.userAuthHandler == nil {
		wsServer.userAuthHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserJoinToRoom(handler func(message UserJoinToRoomMessage, ws WebsocketSender)) {
	if wsServer.userJoinToRoomHandler == nil {
		wsServer.userJoinToRoomHandler = handler
	}
}

func newWsSender(connection *websocket.Conn) *WsSender {
	return &WsSender{connection: connection}
}
