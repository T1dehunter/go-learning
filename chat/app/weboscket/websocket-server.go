package weboscket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketServer struct {
	connections           map[*websocket.Conn]bool
	userConnectHandler    func(message UserConnectMessage, wsSender WebsocketSender)
	userAuthHandler       func(message UserAuthMessage, wsSender WebsocketSender)
	userJoinToRoomHandler func(message UserJoinToRoomMessage, wsSender WebsocketSender)
	userLeaveRoomMessage  func(message UserLeaveRoomMessage, wsSender WebsocketSender)
	userSendRoomMessage   func(message UserSendRoomMessage, wsSender WebsocketSender)
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
	var userConnectMessage UserConnectMessage
	err := json.Unmarshal(messageData, &userConnectMessage)
	if err == nil && userConnectMessage.Name == "user_connect" {
		if wsServer.userConnectHandler != nil {
			sender := NewWsSender(conn)
			wsServer.userConnectHandler(userConnectMessage, sender)
		}
	}

	var userAuthMessage UserAuthMessage
	err = json.Unmarshal(messageData, &userAuthMessage)
	if err == nil && userAuthMessage.Name == "user_auth" {
		if wsServer.userAuthHandler != nil {
			sender := NewWsSender(conn)
			wsServer.userAuthHandler(userAuthMessage, sender)
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

	var userLeaveRoomMessage UserLeaveRoomMessage
	err = json.Unmarshal(messageData, &userLeaveRoomMessage)
	if err == nil && userLeaveRoomMessage.Name == "user_leave_room" {
		if wsServer.userLeaveRoomMessage != nil {
			sender := NewWsSender(conn)
			wsServer.userLeaveRoomMessage(userLeaveRoomMessage, sender)
		}
	}

	var userSendRoomMessage UserSendRoomMessage
	err = json.Unmarshal(messageData, &userSendRoomMessage)
	if err == nil && userSendRoomMessage.Name == "user_send_room_message" {
		if wsServer.userSendRoomMessage != nil {
			sender := NewWsSender(conn)
			wsServer.userSendRoomMessage(userSendRoomMessage, sender)
		}
	}

	if err != nil {
		log.Println("Error unmarshalling message", err)
	}

	//fmt.Println("Received message", message)
}

func (wsServer *WebSocketServer) SubscribeOnUserConnect(handler func(message UserConnectMessage, ws WebsocketSender)) {
	if wsServer.userConnectHandler == nil {
		wsServer.userConnectHandler = handler
	}
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

func (wsServer *WebSocketServer) SubscribeOnUserLeaveRoom(handler func(message UserLeaveRoomMessage, ws WebsocketSender)) {
	if wsServer.userLeaveRoomMessage == nil {
		wsServer.userLeaveRoomMessage = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserSendRoomMessage(handler func(message UserSendRoomMessage, ws WebsocketSender)) {
	if wsServer.userSendRoomMessage == nil {
		wsServer.userSendRoomMessage = handler
	}
}

func newWsSender(connection *websocket.Conn) *WsSender {
	return &WsSender{connection: connection}
}
