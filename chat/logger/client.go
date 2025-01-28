package main

import (
	"chat/logger/console"
	"chat/logger/websocket"
	"io"
)

const (
	stateUserWelcome         = "USER_WELCOME"
	stateUserAuthProcess     = "USER_AUTH_PROCESS"
	stateUserAuthenticated   = "USER_AUTHENTICATED"
	stateUserConnected       = "USER_CONNECTED"
	stateUserJoinedToRoom    = "USER_JOINED_TO_ROOM"
	stateUserSendRoomMessage = "USER_SEND_ROOM_MESSAGE"
	stateUserExit            = "USER_EXIT"
)

type AuthenticatedUser struct {
	ID          int
	Name        string
	AccessToken string
}

type Client struct {
	input     io.Reader
	output    io.Writer
	websocket websocket.Websocket
	wsDataCh  chan string
}

func NewClient(input io.Reader, output io.Writer) *Client {
	return &Client{
		input:  input,
		output: output,
	}
}

func (client *Client) Start() {
	dataChan := make(chan string)
	websocket := websocket.NewWebsocket(dataChan)

	client.websocket = *websocket
	client.websocket.Connect()

	consl := console.NewConsole(dataChan)
	consl.Start()
}
