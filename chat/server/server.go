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
	websocket      *weboscket.WebSocketServer
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
		websocket:      weboscket.NewWebSocketServer(),
	}
}

func (server *Server) Start() {

	http.HandleFunc("/chat", server.websocket.Listen)

	server.websocket.SubscribeOnUserAuth(func(message weboscket.UserAuthMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserAuth(message, ws, server.userService, server.authService, server.roomService)
	})

	server.websocket.SubscribeOnUserConnect(func(message weboscket.UserConnectMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserConnect(message, ws, server.userService, server.authService, server.roomService)
	})

	server.websocket.SubscribeOnUserCreateDirectRoom(func(message weboscket.UserCreateDirectRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserCreateDirectRoom(message, ws, server.userService, server.roomService)
	})

	server.websocket.SubscribeOnUserJoinToRoom(func(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserJoinToRoom(message, ws, server.userService, server.roomService)
	})

	server.websocket.SubscribeOnUserLeaveRoom(func(message weboscket.UserLeaveRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserLeaveRoom(message, ws, server.userService, server.roomService)
	})

	server.websocket.SubscribeOnUserSendDirectMessage(func(message weboscket.UserSendDirectMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserSendDirectMessage(message, ws, server.userService, server.roomService, server.messageService)
	})

	server.websocket.SubscribeOnGetRoomMessages(func(message weboscket.UserGetRoomMessages, ws weboscket.WebsocketSender) {
		handlers.HandleGetRoomMessages(message, ws, server.userService, server.roomService, server.messageService)
	})

	// test handler
	server.websocket.SubscribeOnUserSendRoomMessage(func(message weboscket.UserSendRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserSendRoomMessage(message, ws, server.userService, server.roomService)
	})

	log.Println("Server started!")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
