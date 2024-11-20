package listrooms

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"chat/client/console/types"
	"fmt"
	"strconv"
)

type UserData struct {
	ID    int
	Name  string
	Rooms []types.UserRoom
}
type ListRoomsScreen struct {
	userData        UserData
	renderCh        chan string
	inputTextCh     chan string
	userActionCh    chan interface{}
	userActionResCh chan interface{}
	exitCh          chan interface{}
}

func NewListRoomsScreen(
	renderCh chan string,
	inputTextCh chan string,
	userActionCh chan interface{},
	userActionResCh chan interface{},
) *ListRoomsScreen {
	exitChan := make(chan interface{})
	return &ListRoomsScreen{
		renderCh:        renderCh,
		inputTextCh:     inputTextCh,
		userActionCh:    userActionCh,
		userActionResCh: userActionResCh,
		exitCh:          exitChan,
	}
}

func (screen *ListRoomsScreen) SetUserData(userID int, userName string, rooms []types.UserRoom) {
	screen.userData = UserData{ID: userID, Name: userName, Rooms: rooms}
}

func (screen *ListRoomsScreen) Render() {
	screen.renderContent()
	screen.listenUserInput()
}

func (screen *ListRoomsScreen) renderContent() {
	const template = `
==========================================================================================================
                                 Welcome, %s! 
==========================================================================================================
                             Choose a Room to Start Chatting:
----------------------------------------------------------------------------------------------------------

%s
----------------------------------------------------------------------------------------------------------
Type the room number to enter or 'q' to quit.
`
	rooms := ""
	for _, room := range screen.userData.Rooms {
		rooms += fmt.Sprintf("[%d] %s --- %s \n", room.ID, room.Type, room.Name)
	}

	// fast render
	content := fmt.Sprintf(template, screen.userData.Name, rooms)
	content += symbols.Prompt
	screen.renderCh <- content

	// slow render by character
	//content := fmt.Sprintf(template, screen.userData.Name, rooms)
	//for idx := range content {
	//	currentText += content[:idx + 1]
	//	screen.renderCh <- currentText
	//	time.Sleep(10 * time.Millisecond)
	//}
	//content += symbols.Prompt
	//screen.renderCh <- content
}

func (screen *ListRoomsScreen) Exit() {
	close(screen.exitCh)
}

func (screen *ListRoomsScreen) listenUserInput() {
	go func() {
		for {
			select {
			case text := <-screen.inputTextCh:
				if text == "q" {
					screen.userActionCh <- events.UserChatExit{}
					return
				}
				roomID, err := strconv.Atoi(text)
				if err != nil {
					return
				}
				event := events.UserJoinRoom{RoomID: roomID}
				screen.userActionCh <- event
			case <-screen.exitCh:
				return
			}
		}

	}()
}
