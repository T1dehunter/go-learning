package main

import (
	"chat/logger/console"
	"chat/logger/types"
	"chat/logger/websocket"
	"io"
)

type Client struct {
	input     io.Reader
	output    io.Writer
	websocket *websocket.Websocket
	wsDataCh  chan types.LogEvent
}

func NewClient(input io.Reader, output io.Writer) *Client {
	return &Client{
		input:  input,
		output: output,
	}
}

func (client *Client) Start() {
	dataStream := make(chan types.LogEvent)

	ws := websocket.NewWebsocket(dataStream)
	ws.Connect()

	consl := console.NewConsole(dataStream)
	consl.Start()
}
