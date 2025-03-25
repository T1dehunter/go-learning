package welcome

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"strings"
)

type WelcomeScreen struct {
	renderCh      chan string
	inputTextCh   chan string
	uiActionCh    chan interface{}
	actionResChan chan interface{}
	exitCh        chan interface{}
}

func NewWelcomeScreen(
	renderCh chan string,
	inputTextCh chan string,
	uiActionCh chan interface{},
	actionResChan chan interface{},
) *WelcomeScreen {
	exitChan := make(chan interface{})
	return &WelcomeScreen{
		renderCh:      renderCh,
		inputTextCh:   inputTextCh,
		uiActionCh:    uiActionCh,
		actionResChan: actionResChan,
		exitCh:        exitChan,
	}
}

func (welcomeScreen *WelcomeScreen) Render() {
	welcomeScreen.renderContent()
	welcomeScreen.listenUserInput()
}

func (welcomeScreen *WelcomeScreen) Exit() {
	close(welcomeScreen.exitCh)
}

func (welcomeScreen *WelcomeScreen) renderContent() {
	content := welcomeScreen.buildContent()
	welcomeScreen.renderCh <- content
}

func (welcomeScreen *WelcomeScreen) buildContent() string {
	const logo = `            
                                                __,__
                                       .--.  .-"     "-.  .--.
                                      / .. \/  .-. .-.  \/ .. \
                                      | |  '|  /   Y   \  |'  | |
                                      | \   \  \ 0 | 0 /  /   / |
                                      \ '- ,\.-"""""""-./, -' /
                                       ''-' /_   ^ ^   _\ '-''
                                           |  \._   _./  |
                                           \   \ '~' /   /
                                            '._ '-=-' _.'
                                               '-----'
`
	const text = `
==========================================================================================================
                                Welcome to your first chat, dear user!
==========================================================================================================
We're excited to have you on board! Are you ready to start chatting with your friends? Let's get started!
Do you want to proceed to the authentication step? Type 'yes' to continue or 'no' to exit.             
==========================================================================================================
`
	lines := []string{
		logo,
		text,
		symbols.Prompt,
	}
	content := strings.Join(lines, "\n")
	return content
}

func (welcomeScreen *WelcomeScreen) listenUserInput() {
	go func() {
		for {
			select {
			case text := <-welcomeScreen.inputTextCh:
				confirmation := false
				if text == "yes" {
					confirmation = true
				} else if text == "no" {
					confirmation = false
				}

				event := events.UserChatConfirmed{confirmation}

				welcomeScreen.uiActionCh <- event
			case <-welcomeScreen.exitCh:
				return
			}
		}
	}()
}
