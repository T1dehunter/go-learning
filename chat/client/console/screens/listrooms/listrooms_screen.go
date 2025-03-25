package listrooms

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"chat/client/console/types"
	"fmt"
	"strconv"
	"strings"
)

type UserData struct {
	ID    int
	Name  string
	Rooms []types.Room
}
type ListRoomsScreen struct {
	userData      UserData
	renderCh      chan string
	inputTextCh   chan string
	uiActionCh    chan interface{}
	actionResChan chan interface{}
	exitCh        chan interface{}
}

func NewListRoomsScreen(
	renderCh chan string,
	inputTextCh chan string,
	uiActionCh chan interface{},
	actionResChan chan interface{},
) *ListRoomsScreen {
	exitChan := make(chan interface{})
	return &ListRoomsScreen{
		renderCh:      renderCh,
		inputTextCh:   inputTextCh,
		uiActionCh:    uiActionCh,
		actionResChan: actionResChan,
		exitCh:        exitChan,
	}
}

func (screen *ListRoomsScreen) SetUserData(userID int, userName string, rooms []types.Room) {
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
		var userNames []string
		for _, user := range room.Users {
			userNames = append(userNames, user.Name)
		}
		joinedNames := strings.Join(userNames, ", ")
		rooms += fmt.Sprintf("[%d]%s --- %s | %s \n", room.ID, room.Type, room.Name, joinedNames)
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
					screen.uiActionCh <- events.UserChatExit{}
					return
				}
				roomID, err := strconv.Atoi(text)
				if err != nil {
					return
				}
				event := events.UserJoinRoom{RoomID: roomID}
				screen.uiActionCh <- event
			case <-screen.exitCh:
				return
			}
		}
	}()
}
