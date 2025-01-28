package console

import (
	"bufio"
	"chat/logger/console/types"
	"fmt"
	"os"
	"sync"
	"time"
)

type Console struct {
	logCh     chan string
	logScreen *Screen
}

func NewConsole(logCh chan string) *Console {
	logScreen := NewScreen()
	return &Console{
		logCh:     logCh,
		logScreen: logScreen,
	}
}

func (console *Console) Start() {
	var wg sync.WaitGroup

	console.logScreen.Start()

	console.subscribeOnLogs(&wg)

	time.Sleep(1 * time.Second)

	console.logCh <- "Log 1"
	//time.Sleep(1 * time.Second)
	console.logCh <- "Log 2"
	//time.Sleep(1 * time.Second)
	console.logCh <- "Log 3"
	//console.logCh <- "Log 4"
	// 11 2566 434 dads gfgf fdfdf fd
	//time.Sleep(100 * time.Millisecond)
	//fmt.Println("done")
	wg.Wait()
}

func (console *Console) subscribeOnLogs(wg *sync.WaitGroup) {
	go func() {
		defer wg.Done()
		wg.Add(1)
		for content := range console.logCh {
			//fmt.Println("Log: ", content)
			console.logScreen.Render(types.Event{Title: content, CreatedAt: time.Now().UTC().String()})
			//time.Sleep(100 * time.Millisecond)
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

			//inputText := scanner.Text()

			//console.inputTextCh <- inputText

			//time.Sleep(100 * time.Millisecond)
		}
	}()
}

func (console *Console) print(text string) {
	out := os.Stdout
	fmt.Fprintf(out, "\033[H\033[J")
	fmt.Fprintf(out, text)
}

func (console *Console) DisplayLogScreen() {
	//console.currentScreen.Exit()
	//console.currentScreen = console.authScreen
	//console.currentScreen.Render()
}
