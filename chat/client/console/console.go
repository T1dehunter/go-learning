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
	uiActionCh       chan interface{}
	actionResChan    chan interface{}
	userInputHandler func(message string)
	currentScreen    Screen
}

func NewConsole() *Console {
	renderCh := make(chan string)
	inputTextCh := make(chan string, 1)
	uiActionCh := make(chan interface{})
	actionResChan := make(chan interface{})

	return &Console{
		renderCh:      renderCh,
		inputTextCh:   inputTextCh,
		uiActionCh:    uiActionCh,
		actionResChan: actionResChan,
	}
}

func (console *Console) Start() (chan interface{}, chan interface{}) {
	console.subscribeOnRenderScreen()
	console.subscribeOnInputText()

	console.currentScreen = welcome.NewWelcomeScreen(console.renderCh, console.inputTextCh, console.uiActionCh, console.actionResChan)
	console.currentScreen.Render()

	return console.uiActionCh, console.actionResChan
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

	authScreen := auth.NewAuthScreen(console.renderCh, console.inputTextCh, console.uiActionCh, console.actionResChan)
	authScreen.Render()

	console.currentScreen = authScreen
}

func (console *Console) DisplayListRoomsScreen(userID int, userName string, userRooms []types.Room) {
	console.currentScreen.Exit()

	listRoomsScreen := listrooms.NewListRoomsScreen(console.renderCh, console.inputTextCh, console.uiActionCh, console.actionResChan)
	listRoomsScreen.SetUserData(userID, userName, userRooms)
	listRoomsScreen.Render()

	console.currentScreen = listRoomsScreen
}

func (console *Console) DisplayRoomScreen(userID int, userName string, roomID int, roomName string, users []types.User, messages []types.Message) {
	console.currentScreen.Exit()

	roomScreen := room.NewRoomScreen(console.renderCh, console.inputTextCh, console.uiActionCh, console.actionResChan)
	roomScreen.SetScreenData(room.Data{UserID: userID, RoomID: roomID, RoomName: roomName, Users: users, Messages: messages})
	roomScreen.Render()

	console.currentScreen = roomScreen
}

func (console *Console) DisplayExitScreen() {
	console.currentScreen.Exit()

	exitScreen := exit.NewExitScreen(console.renderCh, console.inputTextCh, console.uiActionCh, console.actionResChan)
	exitScreen.Render()

	console.currentScreen = exitScreen
}

func (console *Console) ReRenderCurrentScreen() {
	console.currentScreen.Render()
}
