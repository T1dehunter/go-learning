package websocket

import (
	"chat/client/websocket/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Websocket struct {
	dataChannel chan string
	connection  *websocket.Conn
	handlers    []func(msg []byte)
	test        string
}

func NewWebsocket(dataChannel chan string) *Websocket {
	return &Websocket{
		dataChannel: dataChannel,
		connection:  nil,
		handlers:    make([]func(msg []byte), 0),
	}
}

func (ws *Websocket) Connect() {
	wsDataChan := ws.initConnect()
	go func() {
		for data := range wsDataChan {
			ws.dataChannel <- data
		}
	}()
}

func (ws *Websocket) initConnect() chan string {
	dataChannel := make(chan string)

	connectToWs := func() {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/chat", nil)

		if err == nil {
			ws.connection = conn
		} else {
			fmt.Println("Error connecting to WebSocket server:", err)
		}
	}

	connectToWs()

	if ws.connection == nil {
		panic("Error on connecting to ws client")
	}

	go func() {
		attempt := 1
		maxAttempts := 5
		for {
			_, wsMessage, err := ws.connection.ReadMessage()
			if err != nil {
				if attempt >= maxAttempts {
					fmt.Println("Max attempts reached, closing connection...")
					close(dataChannel)
					break
				}
				ws.connection.Close()
				time.Sleep(3 * time.Second)
				connectToWs()
				attempt++
			} else {
				for _, handler := range ws.handlers {
					handler(wsMessage)
				}
				dataChannel <- string(wsMessage)
			}
		}
	}()

	return dataChannel
}

func (ws *Websocket) SendUserAuthMsg(name string, password string) *types.UserAuthMsgResponse {
	userAuthMessage := types.UserAuthMsg{
		Type: "user_auth",
		Payload: struct {
			UserName string `json:"userName"`
			Password string `json:"password"`
		}{
			UserName: name,
			Password: password,
		},
	}

	authMessage, err := json.Marshal(userAuthMessage)
	if err != nil {
		log.Fatalf("Error occurred during build auth message. Error: %s", err.Error())
	}
	ws.SendMsg(authMessage)

	var response types.UserAuthMsgResponse
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		fmt.Println("Received message", wsMessage)
		if err != nil || response.Type == "user_authenticated" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SendUserConnectMsg(userID int, accessToken string) *types.UserConnectMsgResponse {
	userConnectMessage := types.UserConnectMsg{
		Type: "user_connect",
		Payload: struct {
			UserID      int    `json:"userID"`
			AccessToken string `json:"accessToken"`
		}{
			UserID:      userID,
			AccessToken: accessToken,
		},
	}

	connectMsg, err := json.Marshal(userConnectMessage)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.SendMsg(connectMsg)

	var response types.UserConnectMsgResponse
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil || response.Type == "user_connected" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SendUserJoinRoomMsg(userID int, roomID int, accessToken string) *types.UserJoinToRoomMsgResponse {
	message := types.UserJoinToRoomMsg{
		Type: "user_join_to_room",
		Payload: struct {
			UserID      int    `json:"userID"`
			RoomID      int    `json:"roomID"`
			AccessToken string `json:"accessToken"`
		}{
			UserID:      userID,
			RoomID:      roomID,
			AccessToken: accessToken,
		},
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.SendMsg(messageJson)

	var response types.UserJoinToRoomMsgResponse
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil {
			panic(err)
		}
		if response.Type == "user_joined_to_room" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SendRoomMsg(userID int, roomID int, text string, accessToken string) *types.UserSendRoomMsgResponse {
	message := types.UserSendRoomMsg{
		Type: "user_send_room_message",
		Payload: types.UserSendRoomMsgPayload{
			UserID:      userID,
			RoomID:      roomID,
			Message:     text,
			AccessToken: accessToken,
		},
	}

	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.SendMsg(messageJson)

	var response types.UserSendRoomMsgResponse
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil {
			panic(err)
		}
		if response.Type == "user_send_room_message" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SubscribeOnRoomMessages(handler func(msg types.UserSendRoomMsgResponse)) {
	ws.handlers = append(ws.handlers, func(msg []byte) {
		var response types.UserSendRoomMsgResponse
		err := json.Unmarshal(msg, &response)
		if err == nil && response.Type == "user_send_room_message" {
			handler(response)
		}
	})
}

func (ws *Websocket) SendUserLeaveRoomMsg(userID int, roomID int) {
	message := types.UserLeaveRoomMsg{
		Type: "user_leave_room",
		Payload: struct {
			UserID int `json:"userID"`
			RoomID int `json:"roomID"`
		}{
			UserID: userID,
			RoomID: roomID,
		},
	}
	messageJson, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.SendMsg(messageJson)
}

func (ws *Websocket) SendLogMsg(message string) {
	logMessage := types.ClientLogMsg{
		Type: "log_message",
		Payload: struct {
			Text string `json:"text"`
		}{
			Text: message,
		},
	}
	messageJson, err := json.Marshal(logMessage)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	ws.SendMsg(messageJson)
}

func (ws *Websocket) SendMsg(message []byte) {
	if ws.connection != nil {
		ws.connection.WriteMessage(websocket.TextMessage, message)
	}
}
