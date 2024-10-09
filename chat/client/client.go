package main

import (
	"chat/client/console"
	"chat/client/userInputParser"
	"chat/client/websocket"
	"encoding/json"
	"fmt"
	"io"
	"log"
)

const (
	USER_AUTHENTICATING = "USER_AUTHENTICATING"
	USER_AUTHENTICATED  = "USER_AUTHENTICATED"
	USER_CONNECTED      = "USER_CONNECTED"
	USER_JOINED_TO_ROOM = "USER_JOINED_TO_ROOM"
)

type UserAuthMessageWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
}

type UserAuthMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		AccessToken string `json:"accessToken"`
	}
}

type UserConnectMessageWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		AccessToken string `json:"accessToken"`
	}
}

type UserConnectMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool `json:"success"`
	}
}

type Client struct {
	userName      string
	currentState  string
	input         io.Reader
	output        io.Writer
	websocket     websocket.Websocket
	wsDataChannel chan string
	wsAuthToken   string
}

func NewClient(input io.Reader, output io.Writer, userName string) *Client {
	wsDataChannel := make(chan string)
	ws := websocket.NewWebsocket(wsDataChannel)

	ws.Connect()

	return &Client{
		userName:      userName,
		currentState:  USER_AUTHENTICATING,
		input:         input,
		output:        output,
		websocket:     *ws,
		wsDataChannel: wsDataChannel,
	}
}

func (client *Client) Start() {

	con := console.NewConsole()

	userInputChannel := con.Start(client.userName)

	userInputParser := userInputParser.NewUserInputParser()

	for userMessage := range userInputChannel {
		if client.isUserAuthenticatingState() {
			userName, password := userInputParser.ParseCredentials(userMessage)
			fmt.Printf("User name: %s, password: %s \n", userName, password)

			client.authenticateUser(userName, password, con)

			fmt.Println("Auth Token after user auth: ", client.wsAuthToken)

			client.connectUser(con)

			con.PrintJoinRoomMessage(client.userName)
		}
		//switch msg := userMessage.(type) {
		//
		//case *console.UserAuthMessage:
		//	fmt.Println("Got user auth message: ", msg)
		//	client.authenticateUser(msg, con)
		//
		//	fmt.Println("Auth Token after user auth: ", client.wsAuthToken)
		//	client.connectUser(con)
		//
		//	con.PrintJoinRoomMessage(client.userName)
		//
		//case *console.UserJoinToRoomMessage:
		//	fmt.Println("Got user join to room message: ", msg)
		//
		//default:
		//	fmt.Println("Unknown user message", userMessage)
		//}
	}
}

func (client *Client) Stop() {
	close(client.wsDataChannel)
}

func (client *Client) authenticateUser(name string, password string, console *console.Console) {
	if !client.isUserAuthenticatingState() {
		return
	}

	//name, password := message.GetPayload()

	// auth:{Sandor Clegane}|{Test1234}
	// auth:{Arya Stark}|{Test4321}

	userAuthMessage := UserAuthMessageWs{
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
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	client.websocket.SendMessage(authMessage)

	var response UserAuthMessageResponseWs
	for wsMessage := range client.wsDataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil {
			panic(err)
		}
		if err == nil && response.Type == "user_authenticated" {
			client.wsAuthToken = response.Payload.AccessToken
			client.setUserAuthenticatedState()
			break
		}
	}
}

func (client *Client) connectUser(console *console.Console) {
	if !client.isUserAuthenticatedState() {
		return
	}

	userConnectMessage := UserConnectMessageWs{
		Type: "user_connect",
		Payload: struct {
			UserID      int    `json:"userID"`
			AccessToken string `json:"accessToken"`
		}{
			UserID:      1,
			AccessToken: client.wsAuthToken,
		},
	}

	authMessage, err := json.Marshal(userConnectMessage)
	if err != nil {
		log.Fatalf("Error occurred during marshaling. Error: %s", err.Error())
	}
	client.websocket.SendMessage(authMessage)

	var response UserConnectMessageResponseWs
	for wsMessage := range client.wsDataChannel {
		err := json.Unmarshal([]byte(wsMessage), &response)
		if err != nil {
			panic(err)
		}
		if err == nil && response.Type == "user_connected" {
			fmt.Printf("User connected: %+v\n", response.Payload)
			client.setUserConnectedState()
			break
		}

	}
}

func (client *Client) isUserAuthenticatingState() bool {
	return client.currentState == USER_AUTHENTICATING
}

func (client *Client) isUserAuthenticatedState() bool {
	return client.currentState == USER_AUTHENTICATED
}

func (client *Client) isUserConnectedState() bool {
	return client.currentState == USER_CONNECTED
}

func (client *Client) isUserJoinedToRoomState() bool {
	return client.currentState == USER_JOINED_TO_ROOM
}

func (client *Client) setUserAuthenticatedState() {
	client.currentState = USER_AUTHENTICATED
}

func (client *Client) setUserConnectedState() {
	client.currentState = USER_CONNECTED
}

func (client *Client) setUserJoinedToRoomState() {
	client.currentState = USER_JOINED_TO_ROOM
}
