package room

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"chat/client/console/types"
	"fmt"
	"time"
)

type UserData struct {
	ID    int
	Name  string
	Rooms []types.UserRoom
}

type RoomScreen struct {
	roomName        string
	renderCh        chan string
	inputTextCh     chan string
	userActionCh    chan interface{}
	userActionResCh chan interface{}
	exitCh          chan interface{}
}

func NewRoomScreen(
	renderCh chan string,
	inputTextCh chan string,
	userActionCh chan interface{},
	userActionResCh chan interface{},
) *RoomScreen {
	exitChan := make(chan interface{})
	return &RoomScreen{
		renderCh:        renderCh,
		inputTextCh:     inputTextCh,
		userActionCh:    userActionCh,
		userActionResCh: userActionResCh,
		exitCh:          exitChan,
	}
}

func (screen *RoomScreen) SetScreenData(roomName string) {
	screen.roomName = roomName
}

func (screen *RoomScreen) Render() {
	//screen.printInitialMessage()
	//
	//screen.listenInput()
	//screen.listenUserActionResult()
	screen.renderContent()
	fmt.Println("Room screen render", screen.roomName)
}

func (screen *RoomScreen) renderContent() {
	const template = `
==========================================================================================================
                            Welcome to room: %s
==========================================================================================================
 Commands:
 [1] Send Message
 [2] View Active Users
 [3] Leave Room

----------------------------------
 Type a command number to proceed.
`
	// slow render by character
	currTittle := ""
	screenText := ""
	for idx := range screen.roomName {
		currTittle = screen.roomName[:idx+1]
		screenText = fmt.Sprintf(template, currTittle)
		screen.renderCh <- screenText
		time.Sleep(50 * time.Millisecond)
	}
	screenText += symbols.Prompt
	screen.renderCh <- screenText

}

func (screen *RoomScreen) Exit() {
	close(screen.exitCh)
}

func (screen *RoomScreen) getLogo() string {
	const AUTH_TEXT = `
==========================================================================================================
                                AUTHORIZATION REQUIRED
==========================================================================================================
`
	return AUTH_TEXT
}

func (screen *RoomScreen) getAskEnterUserName() string {
	return "Enter your username:"
}

func (screen *RoomScreen) printAskEnterPassword() {
	fmt.Println("Enter your password:")
	fmt.Print(symbols.Prompt)
}

func (screen *RoomScreen) listenInput() {
	go func() {
		userName := ""
		password := ""
		for {
			select {
			case text := <-screen.inputTextCh:
				if userName == "" {
					userName = text
					screen.printAskEnterPassword()
				} else if userName != "" && password == "" {
					password = text

					event := events.UserAuthRequest{Username: userName, Password: password}
					screen.userActionCh <- event
				}

			case <-screen.exitCh:
				return
			}
		}

	}()
}

func (screen *RoomScreen) listenUserActionResult() {
	go func() {
		for {
			select {
			case actionResult := <-screen.userActionResCh:
				fmt.Println("room screen get event: ", actionResult)
				//switch event := actionResult.(type) {
				//
				//
				////case events.UserAuthResponse:
				////	fmt.Println("room screen: ", event)
				//}
			}
		}
	}()
}
