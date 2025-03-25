package weboscket

import (
	"github.com/gorilla/websocket"
	"time"
)

type UserId int

type RoomResponse struct {
	users      map[int]*websocket.Conn
	logChannel chan LogMsg
}

func NewRoomResponse(users map[int]*websocket.Conn, logChan chan LogMsg) *RoomResponse {
	return &RoomResponse{users: users, logChannel: logChan}
}

func (response *RoomResponse) SendToAll(message string) {
	data := []byte(message)
	for _, user := range response.users {
		user.WriteMessage(websocket.TextMessage, data)
	}
}

func (response *RoomResponse) SendLog(msg LogMsg) {
	response.logChannel <- msg
}

func (response *RoomResponse) LogText(msg string) {
	response.logChannel <- LogMsg{Title: msg, Type: "server", CreatedAt: time.Now().Format("15:04:05")}
}
