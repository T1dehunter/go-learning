package websocket

import (
	"chat/client/websocket/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Websocket struct {
	dataChannel chan string
	connection  *websocket.Conn
}

func NewWebsocket(dataChannel chan string) *Websocket {
	return &Websocket{
		dataChannel: dataChannel,
		connection:  nil,
	}
}

func (ws *Websocket) Connect() {
	fmt.Println("Console websocket client starting...")
	fmt.Println("\n")

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/chat", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	//defer conn.Close()
	ws.connection = conn

	var message string

	conn.WriteMessage(websocket.TextMessage, []byte(message))

	go func() {
		for {
			_, response, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading from WebSocket:", err)
				ws.dataChannel <- "Error reading from WebSocket"
			} else {
				ws.dataChannel <- string(response)
			}
		}
	}()
}

func (ws *Websocket) Disconnect() {
	if ws.connection != nil {
		ws.connection.Close()
	}
}

func (ws *Websocket) SendUserAuthMessage(name string, password string) *types.UserAuthMessageResponseWs {
	userAuthMessage := types.UserAuthMessageWs{
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
	fmt.Println("Sending auth message", string(authMessage))
	ws.SendMessage(authMessage)

	var response types.UserAuthMessageResponseWs
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		fmt.Println("Received message", wsMessage)
		if err != nil || response.Type == "user_authenticated" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SendUserConnectMessage(userID int, accessToken string) *types.UserConnectMessageResponseWs {
	userConnectMessage := types.UserConnectMessageWs{
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
	ws.SendMessage(connectMsg)

	var response types.UserConnectMessageResponseWs
	for wsMessage := range ws.dataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil || response.Type == "user_connected" {
			break
		}
	}
	return &response
}

func (ws *Websocket) SendUserJoinRoomMessage(userID int, roomID int, accessToken string) *types.UserJoinToRoomMessageResponseWs {
	message := types.UserJoinToRoomMessageWs{
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
	ws.SendMessage(messageJson)

	var response types.UserJoinToRoomMessageResponseWs
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

func (ws *Websocket) SendMessage(message []byte) {
	if ws.connection != nil {
		ws.connection.WriteMessage(websocket.TextMessage, message)
	}
}
