package weboscket

import (
	"github.com/gorilla/websocket"
	"log"
)

type WsSender struct {
	connection     *websocket.Conn
	connectedUsers *map[int]*websocket.Conn
	nameSpace      *map[string][]*websocket.Conn
}

// TODO - strange name, need to be figured out some better
func NewWsSender(conn *websocket.Conn, connectedUsers *map[int]*websocket.Conn, nameSpace *map[string][]*websocket.Conn) *WsSender {
	return &WsSender{connection: conn, connectedUsers: connectedUsers, nameSpace: nameSpace}
}

func (wsSender *WsSender) RegisterConnection(userID int) {
	(*wsSender.connectedUsers)[userID] = wsSender.connection
}

func (wsSender *WsSender) AddUserToNamespace(namespace string) {
	(*wsSender.nameSpace)[namespace] = append((*wsSender.nameSpace)[namespace], wsSender.connection)
}

func (wsSender *WsSender) SendMessageToUser(userID int, message string) {
	//if conn, ok := (*wsSender.connectedUsers)[userID]; ok {
	//	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
	//		log.Println("Error sending message", err)
	//		return
	//	}
	//}
	wsSender.connection.WriteMessage(websocket.TextMessage, []byte(message))
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
