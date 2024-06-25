package main

import (
	"chat/client/console"
	"chat/client/websocket"
	"fmt"
	"io"
)

const (
	AUTH_STATE         = "AUTH_STATE"
	ROOM_CONNECT_STATE = "ROOM_CONNECT_STATE"
)

type Client struct {
	userName     string
	currentState string
	input        io.Reader
	output       io.Writer
	websocket    websocket.Websocket
	dataChannel  chan string
}

func NewClient(input io.Reader, output io.Writer, userName string) *Client {
	dataChannel := make(chan string)
	ws := websocket.NewWebsocket(dataChannel)
	return &Client{
		userName:     userName,
		currentState: AUTH_STATE,
		input:        input,
		output:       output,
		websocket:    *ws,
		dataChannel:  dataChannel,
	}
}

func (client *Client) Start() {
	con := console.NewConsole()

	consoleData := con.Start(client.userName)

	for message := range consoleData {
		switch msg := (*message).(type) {
		case *console.UserAuthMessage:
			fmt.Printf("Got user auth message: %#v\n", msg)
		case *console.UserJoinToRoomMessage:
			fmt.Println("Got user join to room message: ", msg)
		default:
			fmt.Println("Unknown user message", message)
		}

		// auth:{Sandor Clegane}|{Test1234%}

		//if client.currentState == AUTH_STATE && message == "auth" {
		//	client.websocket.Connect()
		//	client.currentState = ROOM_CONNECT_STATE
		//} else if client.currentState == ROOM_CONNECT_STATE {
		//
		//}
		//fmt.Fprintf(client.output, "User input::: %v\n", message)
		//fmt.Fprintf(os.Stdout, ">>> ")
		fmt.Println("User input::: ", *message)
		//message = fmt.Sprintf(`Hello %s`, client.userName)
	}
}
