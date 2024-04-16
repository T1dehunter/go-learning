package weboscket

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsSender struct {
	connection *websocket.Conn
	nameSpace  *map[string][]*websocket.Conn
}

// TODO - strange name, need to be figured out some better
func NewWsSender(conn *websocket.Conn, nameSpace *map[string][]*websocket.Conn) *WsSender {
	return &WsSender{connection: conn, nameSpace: nameSpace}
}

func (wsSender *WsSender) AddUserToNamespace(namespace string) {
	(*wsSender.nameSpace)[namespace] = append((*wsSender.nameSpace)[namespace], wsSender.connection)
}

func (wsSender *WsSender) SendMessageToUser(connectionID int, message string) {
	if err := wsSender.connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("Error sending message", err)
		return
	}
}

func (wsSender *WsSender) SendMessageToRoom(connectionID int, message string) {
	if err := wsSender.connection.Close(); err != nil {
		log.Println("Error closing connection", err)
		return
	}
}

func (wsSender *WsSender) SendMessageToNamespace(namespace string, message string) {
	connections := (*wsSender.nameSpace)[namespace]
	for _, conn := range connections {
		if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
			log.Println("Error sending message", err)
			return
		}
	}
}
