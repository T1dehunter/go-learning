package console

import (
	"bufio"
	"fmt"
	"os"
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

func getAuthMessage(userName string) string {
	return fmt.Sprintf(`Welcome %s!Thank you for using our chat.
Please start by enter chat credentials in following format [name]:[password]`, userName)
	//return fmt.Sprintf(`Hello %s!
	//This is the console client for chat server!
	//Please enter chat credentials in following format ---> auth:{Sandor Clegane}|{Test1234}`, userName)
}

func getJoinRoomMessage(userName string) string {
	return fmt.Sprintf(`Dear %s, please enter the room ID to join in following format ---> join_room:{1}`, userName)
}

type UserMessage interface {
	isMessage()
}

type UserAuthMessage struct {
	name     string
	password string
}

func (userAuthMsg *UserAuthMessage) GetPayload() (string, string) {
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
	dataChannel      chan string
	userInputHandler func(message string)
}

func NewConsole() *Console {
	dataChannel := make(chan string)
	return &Console{
		dataChannel: dataChannel,
	}
}

func (console *Console) Start(userName string) chan string {
	fmt.Println("Console client starting...")
	fmt.Printf(LOGO)
	fmt.Printf(getAuthMessage(userName))
	fmt.Printf("\n")

	in := os.Stdin
	out := os.Stdout
	const PROMPT = ">>> "

	scanner := bufio.NewScanner(in)

	//textParser := NewTextParser()

	fmt.Fprintf(out, PROMPT)

	go func() {
		for {
			scanned := scanner.Scan()

			if !scanned {
				fmt.Println("no scanned")
				close(console.dataChannel)
				return
			}

			inputText := scanner.Text()

			console.dataChannel <- inputText

			//authMessage := textParser.parseAuthMessage(line)
			//if authMessage != nil {
			//	console.dataChannel <- authMessage
			//}

			time.Sleep(100 * time.Millisecond)

			fmt.Fprintf(out, PROMPT)
		}
	}()

	return console.dataChannel
}

func (console *Console) PrintJoinRoomMessage(userName string) {
	fmt.Printf(getJoinRoomMessage(userName))
}

func (console *Console) PrintText(text string) {
	fmt.Println(text)
}
