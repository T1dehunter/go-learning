package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const LOGO = `            __,__
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

type UserMessage interface {
	isMessage()
}

type UserAuthMessage struct {
	name     string
	password string
}

func (userAuthMsg *UserAuthMessage) getPayload() (string, string) {
	return userAuthMsg.name, userAuthMsg.password
}

func (userAuthMsg *UserAuthMessage) isMessage() {}

type UserJoinToRoomMessage struct {
	userID int
	roomID int
}

func (userJoinToRoomMsg *UserJoinToRoomMessage) getPayload() (int, int) {
	return userJoinToRoomMsg.userID, userJoinToRoomMsg.roomID
}

func (userJoinToRoomMsg *UserJoinToRoomMessage) isMessage() {}

type Console struct {
	dataChannel      chan *UserMessage
	userInputHandler func(message string)
}

func NewConsole() *Console {
	dataChannel := make(chan *UserMessage)
	return &Console{
		dataChannel: dataChannel,
	}
}

// returns channel for user messages
func (console *Console) Start(userName string) chan *UserMessage {
	message := fmt.Sprintf(`Hello %s!
This is the console client for chat server!
Please enter chat credentials in following format ---> auth:{username}|{password}`, userName)

	fmt.Println("Console client starting...")
	fmt.Printf(LOGO)
	fmt.Printf(message)
	fmt.Printf("\n")

	//fmt.Printf("Hello %s! This is the console client for chat server!\n", userName)
	//fmt.Printf("Please enter login and password...\n")
	//connectMsg := "{\n \"name\": \"user_connect\", \"payload\": {\"userID\": 1, \"accessToken\": \"Test1234\"}\n}"
	//joinMsg := "{\n \"name\": \"user_join_to_room\", \"payload\": {\"userID\": 1, \"roomID\": 1, \"roomName\": \"Room 1\"}\n}"
	//roomeMessagesMsg := "{\n \"name\": \"user_get_room_messages\", \"payload\": {\"userID\": 1, \"roomID\": 1}\n}"

	in := os.Stdin
	out := os.Stdout
	const PROMPT = ">>> "
	//var userMessage string
	scanner := bufio.NewScanner(in)

	fmt.Fprintf(out, PROMPT)

	go func() {
		for {
			scanned := scanner.Scan()

			if !scanned {
				fmt.Println("no scanned")
				close(console.dataChannel)
				return
			}

			line := scanner.Text()
			userMsg := console.parseUserMessage(line)
			console.dataChannel <- &userMsg

			// sleep for a while
			time.Sleep(100 * time.Millisecond)

			fmt.Fprintf(out, PROMPT)
		}
	}()

	return console.dataChannel
}

func (console *Console) parseUserMessage(message string) UserMessage {
	if strings.Contains(message, "auth") {
		return console.parseUserAuthMessage(message)
	}
	return nil
}

func (console *Console) parseUserAuthMessage(message string) *UserAuthMessage {
	// split message by '|'
	res := strings.Split(message, "|")
	userName, password := res[0], res[1]
	if userName == "" || password == "" {
		fmt.Println("Please enter correct credentials")
		return &UserAuthMessage{
			name:     "",
			password: "",
		}
	}
	return &UserAuthMessage{
		name:     userName,
		password: password,
	}
}

func (console *Console) PrintText(text string) {
	fmt.Printf(text)
}
