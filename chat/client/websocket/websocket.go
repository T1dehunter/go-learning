package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Websocket struct {
	dataChannel chan string
}

func NewWebsocket(dataChannel chan string) *Websocket {
	return &Websocket{
		dataChannel: dataChannel,
	}
}

func (ws *Websocket) Connect() {
	fmt.Println("Console websocket client starting...")
	fmt.Println("\n")

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/chat", nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()

	var message string

	conn.WriteMessage(websocket.TextMessage, []byte(message))

	for {
		_, response, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading from WebSocket:", err)
			ws.dataChannel <- "Error reading from WebSocket"
		} else {
			ws.dataChannel <- string(response)
		}
	}
}
