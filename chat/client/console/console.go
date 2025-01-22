package console

import (
	"bufio"
	"chat/client/console/screens/auth"
	"chat/client/console/screens/exit"
	"chat/client/console/screens/listrooms"
	"chat/client/console/screens/room"
	"chat/client/console/screens/welcome"
	"chat/client/console/types"
	"fmt"
	"os"
)

type PrintFunc func(string)

type Screen interface {
	Render()
	Exit()
}

type UserMessage interface {
	isMessage()
}

type UserAuthMessage struct {
	name     string
	password string
}

func (userAuthMsg *UserAuthMessage) GetPayload() (string, string) {
	return userAuthMsg.name, userAuthMsg.password
}

func (userAuthMsg *UserAuthMessage) isMessage() {}

type UserJoinToRoomMessage struct {
	userID int
	roomID int
}

func (userJoinToRoomMsg *UserJoinToRoomMessage) getPayload() (int, int) {
	return userJoinToRoomMsg.userID, userJoinToRoomMsg.roomID
}

func (userJoinToRoomMsg *UserJoinToRoomMessage) isMessage() {}

type Console struct {
	renderCh         chan string
	inputTextCh      chan string
	userActionCh     chan interface{}
	userActionResCh  chan interface{}
	userInputHandler func(message string)
	welcomeScreen    *welcome.WelcomeScreen
	authScreen       *auth.AuthScreen
	listRoomsScreen  *listrooms.ListRoomsScreen
	roomScreen       *room.RoomScreen
	exitScreen       *exit.ExitScreen
	currentScreen    Screen
}

func NewConsole() *Console {
	renderCh := make(chan string)
	inputTextCh := make(chan string, 1)
	userActionCh := make(chan interface{})
	userActionResCh := make(chan interface{})

	return &Console{
		renderCh:        renderCh,
		inputTextCh:     inputTextCh,
		userActionCh:    userActionCh,
		userActionResCh: userActionResCh,
		welcomeScreen:   welcome.NewWelcomeScreen(renderCh, inputTextCh, userActionCh, userActionResCh),
		authScreen:      auth.NewAuthScreen(renderCh, inputTextCh, userActionCh, userActionResCh),
		listRoomsScreen: listrooms.NewListRoomsScreen(renderCh, inputTextCh, userActionCh, userActionResCh),
		roomScreen:      room.NewRoomScreen(renderCh, inputTextCh, userActionCh, userActionResCh),
		exitScreen:      exit.NewExitScreen(renderCh, inputTextCh, userActionCh, userActionResCh),
	}
}

func (console *Console) Start() (chan interface{}, chan interface{}) {
	console.subscribeOnRenderScreen()
	console.subscribeOnInputText()

	console.currentScreen = console.welcomeScreen
	console.currentScreen.Render()

	return console.userActionCh, console.userActionResCh
}

func (console *Console) subscribeOnRenderScreen() {
	go func() {
		for content := range console.renderCh {
			console.print(content)
		}
	}()
}

func (console *Console) subscribeOnInputText() {
	in := os.Stdin

	//console.print(PROMPT)

	scanner := bufio.NewScanner(in)

	go func() {
		for {
			scanned := scanner.Scan()
			if !scanned {
				//fmt.Println("ttt")
				return
			}

			inputText := scanner.Text()

			console.inputTextCh <- inputText

			//time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (console *Console) print(text string) {
	out := os.Stdout
	fmt.Fprintf(out, "\033[H\033[J")
	fmt.Fprintf(out, text)
}

func (console *Console) DisplayAuthScreen() {
	console.currentScreen.Exit()
	console.currentScreen = console.authScreen
	console.currentScreen.Render()
}

func (console *Console) DisplayListRoomsScreen(userID int, userName string, userRooms []types.Room) {
	console.currentScreen.Exit()

	console.listRoomsScreen.SetUserData(userID, userName, userRooms)

	console.currentScreen = console.listRoomsScreen
	console.currentScreen.Render()
}

func (console *Console) DisplayRoomScreen(userID int, userName string, roomID int, roomName string, users []types.User, messages []types.Message) {
	console.currentScreen.Exit()
	console.roomScreen.SetScreenData(room.Data{UserID: userID, RoomID: roomID, RoomName: roomName, Users: users, Messages: messages})
	console.currentScreen = console.roomScreen
	console.currentScreen.Render()
}

func (console *Console) DisplayExitScreen() {
	console.currentScreen.Exit()
	console.currentScreen = console.exitScreen
	console.currentScreen.Render()
}
