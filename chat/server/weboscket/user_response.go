package weboscket

import (
	"github.com/gorilla/websocket"
	"time"
)

type UserResponse struct {
	conn       *websocket.Conn
	logChannel chan LogMsg
}

func NewUserResponse(conn *websocket.Conn, logChan chan LogMsg) *UserResponse {
	return &UserResponse{conn: conn, logChannel: logChan}
}

func (response *UserResponse) SendMessage(message string) {
	response.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (response *UserResponse) SendMessageToUser(userID int, message string) {
	response.conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (response *UserResponse) SendLog(msg LogMsg) {
	response.logChannel <- msg
}

func (response *UserResponse) LogText(msg string) {
	response.logChannel <- LogMsg{Title: msg, Type: "server", CreatedAt: time.Now().UTC().String()}
}
