package room

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"fmt"
)

const (
	chooseCommand = "CHOOSE_COMMAND"
	enterMessage  = "ENTER_MESSAGE"
)

type Data struct {
	UserID   int
	RoomID   int
	RoomName string
}

type RoomScreen struct {
	state           string
	data            Data
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
		state:           chooseCommand,
		renderCh:        renderCh,
		inputTextCh:     inputTextCh,
		userActionCh:    userActionCh,
		userActionResCh: userActionResCh,
		exitCh:          exitChan,
	}
}

func (screen *RoomScreen) SetScreenData(data Data) {
	screen.data = data
}

func (screen *RoomScreen) Render() {
	screen.renderContent()
	screen.listenUserInput()
	screen.listenUserActionResult()
}

func (screen *RoomScreen) renderContent() {
	const template = `
==========================================================================================================
                            %s
==========================================================================================================
 Commands:
 [1] Send Message
 [2] Leave Room

----------------------------------
 Type a command number to proceed.
`
	fmt.Println("data ::: ", screen.data)
	//
	//const titleTemplate = "Room: #%s"
	//titleText := fmt.Sprintf(titleTemplate, screen.data.roomName)
	//// slow render by character
	//currTittle := ""
	//screenText := ""
	//for idx := range titleText {
	//	currTittle = titleText[:idx+1]
	//	screenText = fmt.Sprintf(template, currTittle)
	//	screen.renderCh <- screenText
	//	time.Sleep(50 * time.Millisecond)
	//}
	//screenText += symbols.Prompt
	//screen.renderCh <- screenText
}

func (screen *RoomScreen) Exit() {
	close(screen.exitCh)
}

func (screen *RoomScreen) setState(nextState string) {
	screen.state = nextState
}

func (screen *RoomScreen) listenUserInput() {
	go func() {
		for {
			select {
			case text := <-screen.inputTextCh:
				if screen.state == chooseCommand && text == "1" {
					screen.printAskEnterMessage()
					screen.setState(enterMessage)
				} else if screen.state == chooseCommand && text == "2" {
					event := events.UserRoomExit{}
					screen.userActionCh <- event
				} else if screen.state == enterMessage && text != "" {
					fmt.Println("room screen send msg: ", text)
					event := events.UserSendMessage{Message: text}
					screen.userActionCh <- event
				}

			case <-screen.exitCh:
				return
			}
		}

	}()
}

func (screen *RoomScreen) printAskEnterMessage() {
	fmt.Println("Enter your message:")
	fmt.Print(symbols.Prompt)
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
