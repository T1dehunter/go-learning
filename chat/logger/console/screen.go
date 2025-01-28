package console

import (
	"chat/logger/console/types"
	"fmt"
)

type UserData struct {
	ID       int
	Username string
}

type Screen struct {
	//logCh chan string
}

func NewScreen(
// logCh chan string,
) *Screen {
	return &Screen{
		//logCh: logCh,
	}
}

func (screen *Screen) Start() {
	const header = `
==========================================================================================================
                                LOGGER SCREEN
==========================================================================================================
`
	fmt.Println(header)
}

func (screen *Screen) Render(event types.Event) {
	text := fmt.Sprintf("Event: %s, CreatedAt: %s", event.Title, event.CreatedAt)
	fmt.Println(text)
	//screen.renderContent()
	//screen.listenLogs()
}

func (screen *Screen) listenLogs() {
	//go func() {
	//	for {
	//		select {
	//		case log := <-screen.logCh:
	//			fmt.Println("Log: ", log)
	//		}
	//	}
	//}()
}
