package main

import (
	"chat/client/console"
	"chat/client/console/events"
	consoleTypes "chat/client/console/types"
	"chat/client/websocket"
	"chat/client/websocket/types"
	"fmt"
	"io"
)

const (
	stateUserWelcome       = "USER_WELCOME"
	stateUserAuthProcess   = "USER_AUTH_PROCESS"
	stateUserAuthenticated = "USER_AUTHENTICATED"
	stateUserConnected     = "USER_CONNECTED"
	stateUserJoinedToRoom  = "USER_JOINED_TO_ROOM"
	stateUserExit          = "USER_EXIT"
)

type AuthenticatedUser struct {
	ID          int
	Name        string
	AccessToken string
}

type Client struct {
	user          *AuthenticatedUser
	userRooms     []consoleTypes.Room
	state         string
	input         io.Reader
	output        io.Writer
	websocket     *websocket.Websocket
	wsDataChannel chan string
}

func NewClient(input io.Reader, output io.Writer) *Client {
	return &Client{
		state:         stateUserWelcome,
		input:         input,
		output:        output,
		wsDataChannel: make(chan string),
	}
}

func (client *Client) Start(user string) {
	ws := websocket.NewWebsocket(client.wsDataChannel)
	ws.Connect()

	client.websocket = ws

	consl := console.NewConsole()

	uiActionCh, actionResChan := consl.Start()

	// TODO - temp for development
	//client.customizeState(consl, user)

	client.listenWsEvents(consl)

	client.listenUserActions(uiActionCh, actionResChan, consl)
}

func (client *Client) Stop() {
	close(client.wsDataChannel)
}

// for development purposes
func (client *Client) customizeState(console *console.Console, userName string) {
	var user *AuthenticatedUser

	if userName == "Sandor" {
		user = &AuthenticatedUser{
			ID:          1,
			Name:        "Sandor Clegane",
			AccessToken: "Test1234",
		}
	}

	if userName == "Arya" {
		user = &AuthenticatedUser{
			ID:          2,
			Name:        "Arya Stark",
			AccessToken: "Test1234",
		}
	}

	if user == nil {
		panic("User not found")
	}

	userRooms := client.connectUser(user)
	client.setUserData(user)
	client.setState(stateUserConnected)
	client.userRooms = userRooms
	console.DisplayListRoomsScreen(user.ID, user.Name, userRooms)
}

func (client *Client) listenWsEvents(consl *console.Console) {
	client.websocket.SubscribeOnRoomMessages(func(msg types.UserSendRoomMsgResponse) {
		if client.state != stateUserJoinedToRoom {
			return
		}

		isCurrUserSender := msg.Payload.UserID == client.user.ID
		if isCurrUserSender {
			return
		}

		var roomUsers []consoleTypes.User
		for _, user := range msg.Payload.Users {
			roomUsers = append(roomUsers, consoleTypes.User{ID: user.ID, Name: user.Name})
		}

		var roomMessages []consoleTypes.Message
		for _, msg := range msg.Payload.Messages {
			message := consoleTypes.Message{
				ID:           msg.ID,
				RoomID:       msg.RoomID,
				CreatorID:    msg.CreatorID,
				CreatorName:  msg.CreatorName,
				ReceiverID:   msg.ReceiverID,
				ReceiverName: msg.ReceiverName,
				Text:         msg.Text,
				CreatedAt:    msg.CreatedAt,
			}
			roomMessages = append(roomMessages, message)
		}
		consl.DisplayRoomScreen(client.user.ID, client.user.Name, msg.Payload.RoomID, msg.Payload.RoomName, roomUsers, roomMessages)
	})
}

func (client *Client) listenUserActions(
	uiActionCh chan interface{},
	actionResChan chan interface{},
	consl *console.Console,
) {
	for userMessage := range uiActionCh {
		switch msg := userMessage.(type) {

		case events.UserChatConfirmed:
			client.handleWelcomeUser(msg, consl)
		case events.UserAuthRequest:
			client.handleAuthUser(msg, consl, actionResChan)
		case events.UserJoinRoom:
			client.handleUserJoinRoom(msg, consl)
		case events.UserRoomExit:
			client.handleUserRoomExit(msg, consl)
		case events.UserSendRoomMessage:
			client.handleUserSendRoomMessage(msg, consl)
		case events.UserChatExit:
			client.handleExitUser(msg, consl)
		default:
			fmt.Println("Unknown user message", userMessage)
		}
	}
}

func (client *Client) handleWelcomeUser(event events.UserChatConfirmed, console *console.Console) {
	if !client.isUserWelcomeState() {
		return
	}
	if event.IsConfirmed {
		client.setState(stateUserAuthProcess)
		console.DisplayAuthScreen()
	} else {
		fmt.Println("User declined chat")
	}
}

func (client *Client) handleAuthUser(msg events.UserAuthRequest, console *console.Console, actionResChan chan interface{}) {
	if !client.isUserAuthProcessState() {
		return
	}
	user := client.authenticateUser(msg.Username, msg.Password, actionResChan)
	if user != nil {
		client.setUserData(user)
		client.setState(stateUserAuthenticated)
		client.handleUserConnect(user, console)
	}
}

func (client *Client) handleUserConnect(user *AuthenticatedUser, console *console.Console) {
	if !client.isUserAuthenticatedState() {
		return
	}
	userRooms := client.connectUser(user)
	client.userRooms = userRooms
	client.setState(stateUserConnected)

	console.DisplayListRoomsScreen(user.ID, user.Name, client.userRooms)
}

func (client *Client) handleUserJoinRoom(event events.UserJoinRoom, console *console.Console) {
	if !client.isUserConnectedState() {
		return
	}

	res := client.websocket.SendUserJoinRoomMsg(client.user.ID, event.RoomID, client.user.AccessToken)

	var roomUsers []consoleTypes.User
	for _, user := range res.Payload.Users {
		roomUsers = append(roomUsers, consoleTypes.User{ID: user.ID, Name: user.Name})
	}

	var roomMessages []consoleTypes.Message
	for _, msg := range res.Payload.Messages {
		message := consoleTypes.Message{
			ID:           msg.ID,
			RoomID:       msg.RoomID,
			CreatorID:    msg.CreatorID,
			CreatorName:  msg.CreatorName,
			ReceiverID:   msg.ReceiverID,
			ReceiverName: msg.ReceiverName,
			Text:         msg.Text,
			CreatedAt:    msg.CreatedAt,
		}
		roomMessages = append(roomMessages, message)
	}

	client.setState(stateUserJoinedToRoom)
	console.DisplayRoomScreen(client.user.ID, client.user.Name, res.Payload.RoomID, res.Payload.RoomName, roomUsers, roomMessages)
}

func (client *Client) handleUserRoomExit(event events.UserRoomExit, console *console.Console) {
	if !client.isUserJoinedToRoomState() {
		return
	}

	client.websocket.SendUserLeaveRoomMsg(client.user.ID, event.RoomID)

	client.setState(stateUserConnected)
	console.DisplayListRoomsScreen(client.user.ID, client.user.Name, client.userRooms)
}

func (client *Client) handleUserSendRoomMessage(event events.UserSendRoomMessage, console *console.Console) {
	if !client.isUserJoinedToRoomState() {
		return
	}

	res := client.websocket.SendRoomMsg(client.user.ID, event.RoomID, event.Message, client.user.AccessToken)

	var roomUsers []consoleTypes.User
	for _, user := range res.Payload.Users {
		roomUsers = append(roomUsers, consoleTypes.User{ID: user.ID, Name: user.Name})
	}

	var roomMessages []consoleTypes.Message
	for _, msg := range res.Payload.Messages {
		message := consoleTypes.Message{
			ID:           msg.ID,
			RoomID:       msg.RoomID,
			CreatorID:    msg.CreatorID,
			CreatorName:  msg.CreatorName,
			ReceiverID:   msg.ReceiverID,
			ReceiverName: msg.ReceiverName,
			Text:         msg.Text,
			CreatedAt:    msg.CreatedAt,
		}
		roomMessages = append(roomMessages, message)
	}

	console.DisplayRoomScreen(
		client.user.ID,
		client.user.Name,
		res.Payload.RoomID,
		res.Payload.RoomName,
		roomUsers,
		roomMessages,
	)
}

func (client *Client) handleExitUser(event events.UserChatExit, console *console.Console) {
	client.setState(stateUserExit)
	console.DisplayExitScreen()
}

func (client *Client) authenticateUser(name string, password string, actionResChan chan interface{}) *AuthenticatedUser {
	response := client.websocket.SendUserAuthMsg(name, password)
	if response == nil || response.Type != "user_authenticated" {
		event := events.UserAuthFailedRes{}
		actionResChan <- event
		return nil
	}

	var user *AuthenticatedUser
	if response.Type == "user_authenticated" {
		user = &AuthenticatedUser{
			ID:          response.Payload.UserID,
			Name:        response.Payload.UserName,
			AccessToken: response.Payload.AccessToken,
		}
	}
	return user
}

func (client *Client) connectUser(user *AuthenticatedUser) []consoleTypes.Room {
	var userRooms []consoleTypes.Room
	response := client.websocket.SendUserConnectMsg(user.ID, user.AccessToken)
	if response == nil || response.Type != "user_connected" {
		return userRooms
	}
	for _, responseRoom := range response.Payload.Rooms {
		users := make([]consoleTypes.User, 0)
		for _, responseUser := range responseRoom.Users {
			users = append(users, consoleTypes.User{ID: responseUser.ID, Name: responseUser.Name})
		}
		userRooms = append(userRooms, consoleTypes.Room{ID: responseRoom.ID, Name: responseRoom.Name, Type: responseRoom.Type, Users: users})
	}
	return userRooms
}

func (client *Client) setUserData(user *AuthenticatedUser) {
	client.user = user
}

func (client *Client) setState(nextState string) {
	client.state = nextState
}

func (client *Client) isUserWelcomeState() bool {
	return client.state == stateUserWelcome
}

func (client *Client) isUserAuthProcessState() bool {
	return client.state == stateUserAuthProcess
}

func (client *Client) isUserAuthenticatedState() bool {
	return client.state == stateUserAuthenticated
}

func (client *Client) isUserConnectedState() bool {
	return client.state == stateUserConnected
}

func (client *Client) isUserJoinedToRoomState() bool {
	return client.state == stateUserJoinedToRoom
}
