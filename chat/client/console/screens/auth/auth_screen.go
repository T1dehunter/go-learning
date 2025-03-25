package auth

import (
	"chat/client/console/events"
	"chat/client/console/symbols"
	"fmt"
	"strings"
	"time"
)

const (
	authenticating   = "AUTHENTICATING"
	reauthenticating = "RE_AUTHENTICATING"
)

type UserData struct {
	ID       int
	Username string
}

type AuthScreen struct {
	state         string
	renderCh      chan string
	inputTextCh   chan string
	uiActionChan  chan interface{}
	actionResChan chan interface{}
	exitCh        chan interface{}
}

func NewAuthScreen(
	renderCh chan string,
	inputTextCh chan string,
	uiActionChan chan interface{},
	actionResChan chan interface{},
) *AuthScreen {
	exitChan := make(chan interface{})
	return &AuthScreen{
		state:         authenticating,
		renderCh:      renderCh,
		inputTextCh:   inputTextCh,
		uiActionChan:  uiActionChan,
		actionResChan: actionResChan,
		exitCh:        exitChan,
	}
}

func (authScreen *AuthScreen) Render() {
	authScreen.renderContent()

	authScreen.listenUserInput()
	authScreen.listenUserActionResult()
}

func (authScreen *AuthScreen) Exit() {
	close(authScreen.exitCh)
}

func (authScreen *AuthScreen) renderContent() {
	const header = `
==========================================================================================================
                                AUTHORIZATION REQUIRED
==========================================================================================================
`
	const askUserName = "Enter your username:"
	rows := []string{
		header,
		askUserName,
		symbols.Prompt,
	}
	content := strings.Join(rows, "\n")

	authScreen.renderCh <- content
}

func (authScreen *AuthScreen) printAskEnterPassword() {
	fmt.Println("Enter your password:")
	fmt.Print(symbols.Prompt)
}

func (authScreen *AuthScreen) listenUserInput() {
	go func() {
		userName := ""
		password := ""
		for {
			select {
			case text := <-authScreen.inputTextCh:
				if userName == "" {
					userName = text
					authScreen.printAskEnterPassword()
				} else if userName != "" && password == "" {
					password = text

					event := events.UserAuthRequest{Username: userName, Password: password}
					authScreen.uiActionChan <- event
				}

				if authScreen.state == reauthenticating {
					if text == "again" {
						userName = ""
						password = ""
						authScreen.state = authenticating
						authScreen.renderContent()
					} else if text == "exit" {
						event := events.UserChatExit{}
						authScreen.uiActionChan <- event
					}
				}

			case <-authScreen.exitCh:
				return
			}
		}

	}()
}

func (authScreen *AuthScreen) listenUserActionResult() {
	go func() {
		for {
			select {
			case actionResult := <-authScreen.actionResChan:
				if _, ok := actionResult.(events.UserAuthFailedRes); ok {
					authScreen.handleUserAuthFailed()
				}
			}
		}
	}()
}

func (authScreen *AuthScreen) handleUserAuthFailed() {
	authScreen.setState(reauthenticating)
	authScreen.renderAuthFailed()
}

func (authScreen *AuthScreen) setState(nextState string) {
	authScreen.state = nextState
}

func (authScreen *AuthScreen) renderAuthFailed() {
	const template = `
==========================================================================================================
                                %s
==========================================================================================================
`
	const title = `Authentication failed`
	currTittle := ""
	var screenText = ""
	for idx := range title {
		currTittle = title[:idx+1]
		screenText = fmt.Sprintf(template, currTittle)
		authScreen.renderCh <- screenText
		time.Sleep(20 * time.Millisecond)
	}
	description := `
We're sorry, but the credentials you provided are incorrect.
Please type 'again' to enter credentials, or type 'exit' to exit chat.`
	for idx := range description {
		currText := description[:idx+1]
		authScreen.renderCh <- screenText + currText
		time.Sleep(20 * time.Millisecond)
	}

	fmt.Println("")
	fmt.Print(symbols.Prompt)
}
