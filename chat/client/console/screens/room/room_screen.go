package room

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"chat/client/console/types"
	"fmt"
	"strings"
	"time"
)

const (
	chooseCommand = "CHOOSE_COMMAND"
	enterMessage  = "ENTER_MESSAGE"
)

type Data struct {
	UserID   int
	RoomID   int
	RoomName string
	Users    []types.User
	Messages []types.Message
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
`
	//fmt.Printf("%+v\n", screen.data)
	titleText := screen.getTitleText()
	currTittle := ""
	screenText := ""
	for idx := range titleText {
		currTittle = titleText[:idx+1]
		if currTittle == "" {
			continue
		}
		screenText = fmt.Sprintf(template, currTittle)
		screen.renderCh <- screenText
		time.Sleep(20 * time.Millisecond)
	}

	messagesText := screen.getMessagesText()

	screenText += messagesText

	screenText += `
----------------------------------------------------------------------------------------------------------
Commands:
[1] Type your message
[2] Exit Room
----------------------------------------------------------------------------------------------------------
`

	screenText += symbols.Prompt

	screen.renderCh <- screenText
}

func (screen *RoomScreen) getTitleText() string {
	var usersNames []string
	for _, user := range screen.data.Users {
		usersNames = append(usersNames, "@"+user.Name)
	}
	formattedNames := strings.Join(usersNames, ", ")
	const titleTemplate = "Room: #{roomName}, Users: {usersNames}"
	replacer := strings.NewReplacer("{roomName}", screen.data.RoomName, "{usersNames}", formattedNames)
	titleText := replacer.Replace(titleTemplate)
	return titleText
}

func (screen *RoomScreen) getMessagesText() string {
	var messagesText string
	for _, message := range screen.data.Messages {
		messagesText += fmt.Sprintf("User: @%s [%s]: %s\n", message.CreatorName, message.CreatedAt, message.Text)
	}
	return messagesText
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
					//fmt.Println("room screen send msg: ", text)
					event := events.UserSendRoomMessage{RoomID: screen.data.RoomID, Message: text}
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
