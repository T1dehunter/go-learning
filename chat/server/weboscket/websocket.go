package weboscket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type WebSocket struct {
	users                        map[int]*websocket.Conn
	userConnectHandler           func(message UserConnectMsg, response *Response)
	userAuthHandler              func(message UserAuthMsg, response *Response)
	userCreateDirectRoomHandler  func(message UserCreateRoomMsg, response *Response)
	userJoinToRoomHandler        func(message UserJoinToRoomMsg, response *Response)
	userSendRoomMessageHandler   func(message UserSendRoomMsg, response *RoomResponse)
	userLeaveRoomHandler         func(message UserLeaveRoomMsg, response *Response)
	userSendDirectMessageHandler func(message UserSendDirectMsg, response *Response)
	userGetRoomMessagesHandler   func(message UserGetListRoomMsg, response *Response)
	clientLogMsgHandler          func(message ClientLogMsg, response *Response)
	logsChannel                  chan LogMsg
}

func NewWebSocket() *WebSocket {
	users := make(map[int]*websocket.Conn)
	logsChannel := make(chan LogMsg, 100)

	return &WebSocket{users: users, logsChannel: logsChannel}
}

func (webSocket *WebSocket) HandleEvents(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("conn upgrade error: ", err)
		return
	}

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("server ws error: %v", err)
			} else {
				log.Println("some error: ", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		webSocket.HandleMessage(conn, message)
	}
}

func (webSocket *WebSocket) StreamLogs(w http.ResponseWriter, r *http.Request) {
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

	go func() {
		time.Sleep(1 * time.Second)
		currentTime := time.Now()
		formattedTime := currentTime.Format("15:04:05")
		webSocket.logsChannel <- LogMsg{Title: "Test message 1", Type: "server", CreatedAt: formattedTime}
		webSocket.logsChannel <- LogMsg{Title: "Test message 2", Type: "server", CreatedAt: formattedTime}
		webSocket.logsChannel <- LogMsg{Title: "Test message 3", Type: "server", CreatedAt: formattedTime}
		webSocket.logsChannel <- LogMsg{Title: "Test message 4", Type: "server", CreatedAt: formattedTime}
		webSocket.logsChannel <- LogMsg{Title: "Test message 5", Type: "server", CreatedAt: formattedTime}
	}()

	for logEvent := range webSocket.logsChannel {
		logEventJson, _ := json.Marshal(logEvent)

		fmt.Println("send log message", string(logEventJson))

		w, err := conn.NextWriter(1)
		if err != nil {
			fmt.Println("Error getting writer", err)
			return
		}

		w.Write(logEventJson)
	}
}

func (webSocket *WebSocket) HandleMessage(conn *websocket.Conn, messageData []byte) {
	fmt.Println("conn.LocalAddr", conn.LocalAddr())

	var userAuthMessage UserAuthMsg
	err := json.Unmarshal(messageData, &userAuthMessage)
	if err == nil && userAuthMessage.Type == "user_auth" {
		if webSocket.userAuthHandler != nil {
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userAuthHandler(userAuthMessage, response)
		}
	}

	var userConnectMessage UserConnectMsg
	err = json.Unmarshal(messageData, &userConnectMessage)
	if err == nil && userConnectMessage.Type == "user_connect" {
		webSocket.users[userConnectMessage.Payload.UserID] = conn
		fmt.Println("USERS", webSocket.users)
		if webSocket.userConnectHandler != nil {
			webSocket.log("User connected")
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userConnectHandler(userConnectMessage, response)
		}
	}

	var userCreateDirectRoomMessage UserCreateRoomMsg
	err = json.Unmarshal(messageData, &userCreateDirectRoomMessage)
	if err == nil && userCreateDirectRoomMessage.Type == "user_create_direct_room" {
		if webSocket.userCreateDirectRoomHandler != nil {
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userCreateDirectRoomHandler(userCreateDirectRoomMessage, response)
		}

	}

	var userJoinToRoomMessage UserJoinToRoomMsg
	err = json.Unmarshal(messageData, &userJoinToRoomMessage)
	if err == nil && userJoinToRoomMessage.Type == "user_join_to_room" {
		if webSocket.userJoinToRoomHandler != nil {
			webSocket.log("User join to room: " + string(messageData))
			// there we need save connections by room, in order to send messages to all users in room
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userJoinToRoomHandler(userJoinToRoomMessage, response)
		}
	}

	var userSendRoomMessage UserSendRoomMsg
	err = json.Unmarshal(messageData, &userSendRoomMessage)
	if err == nil && userSendRoomMessage.Type == "user_send_room_message" {
		if webSocket.userSendRoomMessageHandler != nil {
			response := NewRoomResponse(webSocket.users, webSocket.logsChannel)
			//response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userSendRoomMessageHandler(userSendRoomMessage, response)
		}
	}

	var userLeaveRoomMessage UserLeaveRoomMsg
	err = json.Unmarshal(messageData, &userLeaveRoomMessage)
	if err == nil && userLeaveRoomMessage.Type == "user_leave_room" {
		if webSocket.userLeaveRoomHandler != nil {
			webSocket.log("User leave room: " + string(messageData))
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userLeaveRoomHandler(userLeaveRoomMessage, response)
		}
	}

	var userSendDirectMessage UserSendDirectMsg
	err = json.Unmarshal(messageData, &userSendDirectMessage)
	if err == nil && userSendDirectMessage.Type == "direct_message" {
		if webSocket.userSendDirectMessageHandler != nil {
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userSendDirectMessageHandler(userSendDirectMessage, response)
		}
	}

	var userGetRoomMessages UserGetListRoomMsg
	err = json.Unmarshal(messageData, &userGetRoomMessages)
	if err == nil && userGetRoomMessages.Type == "user_get_room_messages" {
		if webSocket.userGetRoomMessagesHandler != nil {
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.userGetRoomMessagesHandler(userGetRoomMessages, response)
		}
	}

	var testMessage ClientLogMsg
	err = json.Unmarshal(messageData, &testMessage)
	if err == nil && testMessage.Type == "log_message" {
		if webSocket.clientLogMsgHandler != nil {
			response := NewResponse(conn, webSocket.logsChannel)
			webSocket.clientLogMsgHandler(testMessage, response)
		}
	}

	if err != nil {
		log.Println("Error unmarshalling message", err)
	}
}

func (webSocket *WebSocket) SubscribeOnUserConnect(handler func(message UserConnectMsg, response *Response)) {
	if webSocket.userConnectHandler == nil {
		webSocket.userConnectHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserAuth(handler func(message UserAuthMsg, response *Response)) {
	if webSocket.userAuthHandler == nil {
		webSocket.userAuthHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserCreateDirectRoom(handler func(message UserCreateRoomMsg, response *Response)) {
	if webSocket.userCreateDirectRoomHandler == nil {
		webSocket.userCreateDirectRoomHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserJoinToRoom(handler func(message UserJoinToRoomMsg, response *Response)) {
	if webSocket.userJoinToRoomHandler == nil {
		webSocket.userJoinToRoomHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserSendRoomMessage(handler func(message UserSendRoomMsg, response *RoomResponse)) {
	if webSocket.userSendRoomMessageHandler == nil {
		webSocket.userSendRoomMessageHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserLeaveRoom(handler func(message UserLeaveRoomMsg, response *Response)) {
	if webSocket.userLeaveRoomHandler == nil {
		webSocket.userLeaveRoomHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnUserSendDirectMessage(handler func(message UserSendDirectMsg, response *Response)) {
	if webSocket.userSendDirectMessageHandler == nil {
		webSocket.userSendDirectMessageHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnGetRoomMessages(handler func(message UserGetListRoomMsg, response *Response)) {
	if webSocket.userGetRoomMessagesHandler == nil {
		webSocket.userGetRoomMessagesHandler = handler
	}
}

func (webSocket *WebSocket) SubscribeOnClientLogMsg(handler func(message ClientLogMsg, response *Response)) {
	if webSocket.clientLogMsgHandler == nil {
		webSocket.clientLogMsgHandler = handler
	}
}

func (webSocket *WebSocket) SendMessageToUser(userID int, message string) {
	//if conn, ok := (*wsSender.connectedUsers)[userID]; ok {
	//	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
	//		log.Println("Error sending message", err)
	//		return
	//	}
	//}
	//webSocket.connection.WriteMessage(websocket.TextMessage, []byte(message))
}

func (webSocket *WebSocket) log(text string) {
	currentTime := time.Now()
	formattedTime := currentTime.Format("15:04:05")
	//fmt.Println("SEND SEND SEND SEND SEND SEND SEND SEND SEND SEND SEND LOG:", text, time.Now().UTC().String())
	webSocket.logsChannel <- LogMsg{Title: text, Type: "server", CreatedAt: formattedTime}
}
