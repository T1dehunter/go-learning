package console

import (
	"chat/logger/types"
	"fmt"
)

type Screen struct {
	//logCh chan string
}

func NewScreen() *Screen {
	return &Screen{}
}

func (screen *Screen) Start() {
	const header = `
==========================================================================================================
                                LOGGER SCREEN
==========================================================================================================
`
	fmt.Println(header)
}

func (screen *Screen) Render(event types.LogEvent) {
	//text := fmt.Sprintf("Type: %s, Title: %s, CreatedAt: %s", event.Type, event.Title, event.CreatedAt)
	eventType := "CLIENT"
	if event.Type == types.ServerEvent {
		eventType = "SERVER"
	}
	text := fmt.Sprintf("[%s] [%s] %s", event.CreatedAt, eventType, event.Title)

	fmt.Println(text)
}
