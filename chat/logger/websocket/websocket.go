package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
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

func (ws *Websocket) initConnect() error {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/chat", nil)
	if err == nil {
		ws.connection = conn
	}
	return err
}

func (ws *Websocket) Connect() {
	fmt.Println("Console websocket client starting...")
	fmt.Println("\n")

	err := ws.initConnect()
	if err != nil {
		message := fmt.Sprintf("Error connecting to WebSocket server: %s", err.Error())
		fmt.Println(message)
		ws.reConnect()
	}

	conn := ws.connection

	var message string

	conn.WriteMessage(websocket.TextMessage, []byte(message))

	go func() {
		for {
			_, response, err := conn.ReadMessage()
			if err != nil {
				log.Println("Error reading from WebSocket:", err)
				ws.dataChannel <- "Error reading from WebSocket"
				//ws.reConnect()
			} else {
				ws.dataChannel <- string(response)
			}
		}
	}()
}

func (ws *Websocket) reConnect() {
	fmt.Println("Start reconnecting to the WebSocket server ...")

	maxAttempts := 30
	currentAttempt := 0

	for currentAttempt < maxAttempts {
		message := fmt.Sprintf("Reconnecting to the WebSocket server, attempt: %d", currentAttempt)
		fmt.Println(message)

		err := ws.initConnect()
		if err == nil {
			break
		}

		currentAttempt++
		time.Sleep(1 * time.Second)
	}
}

func (ws *Websocket) Disconnect() {
	if ws.connection != nil {
		ws.connection.Close()
	}
}
