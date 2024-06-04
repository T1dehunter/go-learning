package weboscket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type WebSocketServer struct {
	connections                  map[*websocket.Conn]bool
	connectedUsers               map[int]*websocket.Conn
	nameSpace                    map[string][]*websocket.Conn
	userConnectHandler           func(message UserConnectMessage, wsSender WebsocketSender)
	userAuthHandler              func(message UserAuthMessage, wsSender WebsocketSender)
	userCreateDirectRoomHandler  func(message UserCreateDirectRoomMessage, wsSender WebsocketSender)
	userJoinToRoomHandler        func(message UserJoinToRoomMessage, wsSender WebsocketSender)
	userLeaveRoomHandler         func(message UserLeaveRoomMessage, wsSender WebsocketSender)
	userSendDirectMessageHandler func(message UserSendDirectMessage, wsSender WebsocketSender)
	userSendRoomMessageHandler   func(message UserSendRoomMessage, wsSender WebsocketSender)
	UserGetRoomMessagesHandler   func(message UserGetRoomMessages, wsSender WebsocketSender)
}

func NewWebSocketServer() *WebSocketServer {
	connections := make(map[*websocket.Conn]bool)
	connectedUsers := make(map[int]*websocket.Conn)
	nameSpace := make(map[string][]*websocket.Conn)
	return &WebSocketServer{connections: connections, connectedUsers: connectedUsers, nameSpace: nameSpace}
}

func (wsServer *WebSocketServer) Listen(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
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
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userConnectHandler(userConnectMessage, sender)
		}
	}

	var userAuthMessage UserAuthMessage
	err = json.Unmarshal(messageData, &userAuthMessage)
	if err == nil && userAuthMessage.Name == "user_auth" {
		if wsServer.userAuthHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userAuthHandler(userAuthMessage, sender)
		}
	}

	var userCreateDirectRoomMessage UserCreateDirectRoomMessage
	err = json.Unmarshal(messageData, &userCreateDirectRoomMessage)
	if err == nil && userCreateDirectRoomMessage.Name == "user_create_direct_room" {
		if wsServer.userCreateDirectRoomHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userCreateDirectRoomHandler(userCreateDirectRoomMessage, sender)
		}
	}

	var userJoinToRoomMessage UserJoinToRoomMessage
	err = json.Unmarshal(messageData, &userJoinToRoomMessage)
	if err == nil && userJoinToRoomMessage.Name == "user_join_to_room" {
		if wsServer.userJoinToRoomHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userJoinToRoomHandler(userJoinToRoomMessage, sender)
		}
	}

	var userLeaveRoomMessage UserLeaveRoomMessage
	err = json.Unmarshal(messageData, &userLeaveRoomMessage)
	if err == nil && userLeaveRoomMessage.Name == "user_leave_room" {
		if wsServer.userLeaveRoomHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userLeaveRoomHandler(userLeaveRoomMessage, sender)
		}
	}

	var userSendDirectMessage UserSendDirectMessage
	err = json.Unmarshal(messageData, &userSendDirectMessage)
	if err == nil && userSendDirectMessage.Name == "direct_message" {
		if wsServer.userSendDirectMessageHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userSendDirectMessageHandler(userSendDirectMessage, sender)
		}
	}

	var userSendRoomMessage UserSendRoomMessage
	err = json.Unmarshal(messageData, &userSendRoomMessage)
	if err == nil && userSendRoomMessage.Name == "room_message" {
		if wsServer.userSendRoomMessageHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.userSendRoomMessageHandler(userSendRoomMessage, sender)
		}
	}

	var userGetRoomMessages UserGetRoomMessages
	err = json.Unmarshal(messageData, &userGetRoomMessages)
	if err == nil && userGetRoomMessages.Name == "user_get_room_messages" {
		if wsServer.UserGetRoomMessagesHandler != nil {
			sender := NewWsSender(conn, &wsServer.connectedUsers, &wsServer.nameSpace)
			wsServer.UserGetRoomMessagesHandler(userGetRoomMessages, sender)
		}
	}

	if err != nil {
		log.Println("Error unmarshalling message", err)
	}

	fmt.Println("NameSpace", wsServer.nameSpace)
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

func (wsServer *WebSocketServer) SubscribeOnUserCreateDirectRoom(handler func(message UserCreateDirectRoomMessage, ws WebsocketSender)) {
	if wsServer.userCreateDirectRoomHandler == nil {
		wsServer.userCreateDirectRoomHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserJoinToRoom(handler func(message UserJoinToRoomMessage, ws WebsocketSender)) {
	if wsServer.userJoinToRoomHandler == nil {
		wsServer.userJoinToRoomHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserLeaveRoom(handler func(message UserLeaveRoomMessage, ws WebsocketSender)) {
	if wsServer.userLeaveRoomHandler == nil {
		wsServer.userLeaveRoomHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserSendDirectMessage(handler func(message UserSendDirectMessage, ws WebsocketSender)) {
	if wsServer.userSendDirectMessageHandler == nil {
		wsServer.userSendDirectMessageHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnUserSendRoomMessage(handler func(message UserSendRoomMessage, ws WebsocketSender)) {
	if wsServer.userSendRoomMessageHandler == nil {
		wsServer.userSendRoomMessageHandler = handler
	}
}

func (wsServer *WebSocketServer) SubscribeOnGetRoomMessages(handler func(message UserGetRoomMessages, ws WebsocketSender)) {
	if wsServer.UserGetRoomMessagesHandler == nil {
		wsServer.UserGetRoomMessagesHandler = handler
	}
}

func newWsSender(connection *websocket.Conn) *WsSender {
	return &WsSender{connection: connection}
}
