package websocket

import (
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

func (ws *Websocket) SendMessage(message []byte) {
	if ws.connection != nil {
		ws.connection.WriteMessage(websocket.TextMessage, message)
	}
}
