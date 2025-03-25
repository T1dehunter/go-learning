package main

import (
	"chat/server/components/auth"
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/database"
	"chat/server/handlers"
	"chat/server/weboscket"
	"log"
	"net/http"
)

type Server struct {
	userService    *user.UserService
	authService    *auth.AuthService
	roomService    *room.RoomService
	messageService *message.MessageService
	websocket      *weboscket.WebSocket
}

func NewServer() *Server {
	dbClient := database.Connect()
	userRepository := user.NewUserRepository(dbClient)
	roomRepository := room.NewRoomRepository(dbClient)
	messageRepository := message.NewMessageRepository(dbClient)

	return &Server{
		authService:    auth.NewAuthService(userRepository),
		userService:    user.NewUserService(userRepository),
		roomService:    room.NewRoomService(roomRepository),
		messageService: message.NewMessageService(messageRepository),
		websocket:      weboscket.NewWebSocket(),
	}
}

func (server *Server) Start() {
	http.HandleFunc("/chat", server.websocket.HandleEvents)

	http.HandleFunc("/logs", server.websocket.StreamLogs)

	server.websocket.SubscribeOnUserAuth(func(message weboscket.UserAuthMsg, response *weboscket.Response) {
		handlers.HandleUserAuth(message, server.userService, response)
	})

	server.websocket.SubscribeOnUserConnect(func(message weboscket.UserConnectMsg, response *weboscket.Response) {
		handlers.HandleUserConnect(message, server.userService, server.authService, server.roomService, response)
	})

	server.websocket.SubscribeOnUserCreateDirectRoom(func(message weboscket.UserCreateRoomMsg, response *weboscket.Response) {
		handlers.HandleUserCreateRoom(message, server.userService, server.roomService, response)
	})

	server.websocket.SubscribeOnUserJoinToRoom(func(message weboscket.UserJoinToRoomMsg, response *weboscket.Response) {
		handlers.HandleUserJoinRoom(message, server.userService, server.roomService, server.messageService, response)
	})

	server.websocket.SubscribeOnUserSendRoomMessage(func(message weboscket.UserSendRoomMsg, response *weboscket.RoomResponse) {
		handlers.HandleUserSendRoomMessage(message, server.userService, server.roomService, server.messageService, response)
	})

	server.websocket.SubscribeOnUserLeaveRoom(func(message weboscket.UserLeaveRoomMsg, response *weboscket.Response) {
		handlers.HandleUserLeaveRoom(message, server.userService, server.roomService, response)
	})

	server.websocket.SubscribeOnUserSendDirectMessage(func(message weboscket.UserSendDirectMsg, response *weboscket.Response) {
		handlers.HandleUserSendDirectMessage(message, server.userService, server.roomService, server.messageService, response)
	})

	server.websocket.SubscribeOnGetRoomMessages(func(message weboscket.UserGetListRoomMsg, response *weboscket.Response) {
		handlers.HandleGetRoomMessages(message, server.userService, server.roomService, server.messageService, response)
	})

	server.websocket.SubscribeOnClientLogMsg(func(message weboscket.ClientLogMsg, response *weboscket.Response) {
		handlers.HandleClientLogMsg(message, response)
	})

	log.Println("websocket server started at :3000")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
