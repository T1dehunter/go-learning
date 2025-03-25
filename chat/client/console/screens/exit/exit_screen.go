package exit

import (
	"chat/client/console/symbols"
	"fmt"
	"os"
	"time"
)

type ExitScreen struct {
	renderCh      chan string
	inputTextCh   chan string
	uiActionCh    chan interface{}
	actionResChan chan interface{}
	exit          chan interface{}
}

func NewExitScreen(renderCh chan string, inputTextCh chan string, uiActionCh chan interface{}, actionResChan chan interface{}) *ExitScreen {
	exitChan := make(chan interface{})
	return &ExitScreen{
		renderCh:      renderCh,
		inputTextCh:   inputTextCh,
		uiActionCh:    uiActionCh,
		actionResChan: actionResChan,
		exit:          exitChan,
	}
}

func (exitScreen *ExitScreen) Render() {
	exitScreen.renderContent()

	exitScreen.listenUserInput()
}

func (exitScreen *ExitScreen) renderContent() {
	const template = `
==========================================================================================================
                                %s
==========================================================================================================
`
	const title = `Thank You for Using Our Chat!`
	currTittle := ""
	screenText := ""
	for idx := range title {
		currTittle = title[:idx+1]
		screenText = fmt.Sprintf(template, currTittle)
		exitScreen.renderCh <- screenText
		time.Sleep(20 * time.Millisecond)
	}

	description := `
We hope you had a great experience.

→ Remember, you can always come back to chat and reconnect!  
→ Have feedback? Let us know next time.
Press *Enter* to exit.`
	for idx := range description {
		currText := description[:idx+1]
		exitScreen.renderCh <- screenText + currText
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Println("")
	fmt.Print(symbols.Prompt)
}

func (exitScreen *ExitScreen) listenUserInput() {
	go func() {
		for {
			select {
			case text := <-exitScreen.inputTextCh:
				if text == "" {
					exitScreen.Exit()
				}
			case <-exitScreen.exit:
				return
			}
		}
	}()
}

func (exitScreen *ExitScreen) Exit() {
	fmt.Println("")
	fmt.Println("Exiting the chat...")
	time.Sleep(1 * time.Second)
	os.Exit(0)
}
