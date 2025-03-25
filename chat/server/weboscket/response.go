package weboscket

import (
	"github.com/gorilla/websocket"
	"time"
)

type Response struct {
	userInitiatorID  int
	userInitiatorCon *websocket.Conn
	logChannel       chan LogMsg
}

func NewResponse(conn *websocket.Conn, logChan chan LogMsg) *Response {
	return &Response{userInitiatorCon: conn, logChannel: logChan}
}

func (response *Response) SendMessage(message string) {
	response.userInitiatorCon.WriteMessage(websocket.TextMessage, []byte(message))
}

func (response *Response) SendMessageToUser(userID int, message string) {
	response.userInitiatorCon.WriteMessage(websocket.TextMessage, []byte(message))
}

func (response *Response) SendLog(msg LogMsg) {
	response.logChannel <- msg
}

func (response *Response) LogText(msg string) {
	response.logChannel <- LogMsg{Title: msg, Type: "server", CreatedAt: time.Now().Format("15:04:05")}
}
