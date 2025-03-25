package console

import (
	"chat/logger/types"
	"sync"
)

type Console struct {
	dataStream chan types.LogEvent
	logScreen  *Screen
}

func NewConsole(dataStream chan types.LogEvent) *Console {
	logScreen := NewScreen()
	return &Console{
		dataStream: dataStream,
		logScreen:  logScreen,
	}
}

func (console *Console) Start() {
	var wg sync.WaitGroup

	console.logScreen.Start()

	wg.Add(1)

	console.subscribeOnLogs(&wg)

	wg.Wait()
}

func (console *Console) subscribeOnLogs(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		for logEvent := range console.dataStream {
			console.logScreen.Render(logEvent)
		}
	}()
}
