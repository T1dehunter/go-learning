package websocket

import (
	"chat/logger/types"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type WsData = []byte

type Websocket struct {
	streamChan chan types.LogEvent
}

func NewWebsocket(streamChan chan types.LogEvent) *Websocket {
	return &Websocket{
		streamChan: streamChan,
	}
}

func (ws *Websocket) Connect() {
	wsDataChan := ws.initConnect()

	go func() {
		for data := range wsDataChan {
			if data == nil {
				continue
			}
			var response types.LogEvent
			err := json.Unmarshal(data, &response)
			if err != nil {
				log.Println("error on parse message from websocket:", err)
			} else {
				ws.streamChan <- response
			}
		}
	}()
}

func (ws *Websocket) initConnect() chan WsData {
	var connection *websocket.Conn

	defer func() {
		fmt.Println("closing connection...")
		if connection != nil {
			connection.Close()
		}
	}()

	dataChannel := make(chan WsData)

	connectToWs := func() {
		conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:3000/logs", nil)

		if err == nil {
			fmt.Println("logger client connected to websocket server")
			connection = conn
		} else {
			fmt.Println("Error connecting to WebSocket server:", err)
		}

	}

	connectToWs()

	if connection == nil {
		panic("Error connecting to WebSocket server")
	}

	go func() {
		attempt := 1
		maxAttempts := 5
		for {
			_, wsMessage, err := connection.ReadMessage()
			if err != nil {
				if attempt >= maxAttempts {
					close(dataChannel)
					break
				}
				connection.Close()

				time.Sleep(2 * time.Second)

				connectToWs()
				attempt++
			} else {
				dataChannel <- wsMessage
			}
		}
	}()

	return dataChannel
}
