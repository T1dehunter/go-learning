package app

import (
	"chat/app/components/auth"
	"chat/app/components/message"
	"chat/app/components/room"
	"chat/app/components/user"
	"chat/app/database"
	"chat/app/handlers"
	"chat/app/weboscket"
	"log"
	"net/http"
)

//import "chat/app/components/auth"
//import "chat/app/components/room"
//import "chat/app/weboscket"

//import "chat/app/handlers"

type App struct {
	userService    *user.UserService
	authService    *auth.AuthService
	roomService    *room.RoomService
	messageService *message.MessageService
	wsServer       *weboscket.WebSocketServer
}

func NewApp() *App {
	client := database.Connect()

	database.TestFind()

	userRepository := user.NewUserRepository(client)
	roomRepository := room.NewRoomRepository(client)
	messageRepository := message.NewMessageRepository(client)

	return &App{
		authService:    auth.NewAuthService(userRepository),
		userService:    user.NewUserService(userRepository),
		roomService:    room.NewRoomService(roomRepository),
		messageService: message.NewMessageService(messageRepository),
		wsServer:       weboscket.NewWebSocketServer(),
	}
}

func (app *App) Start() {

	http.HandleFunc("/chat", app.wsServer.Listen)

	app.wsServer.SubscribeOnUserConnect(func(message weboscket.UserConnectMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserConnect(message, ws, app.userService, app.authService)
	})

	app.wsServer.SubscribeOnUserAuth(func(message weboscket.UserAuthMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserAuth(message, ws, app.userService, app.authService, app.roomService)
	})

	app.wsServer.SubscribeOnUserCreateDirectRoom(func(message weboscket.UserCreateDirectRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserCreateDirectRoom(message, ws, app.userService, app.roomService)
	})

	app.wsServer.SubscribeOnUserJoinToRoom(func(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserJoinToRoom(message, ws, app.userService, app.roomService)
	})

	app.wsServer.SubscribeOnUserLeaveRoom(func(message weboscket.UserLeaveRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserLeaveRoom(message, ws, app.userService, app.roomService)
	})

	app.wsServer.SubscribeOnUserSendDirectMessage(func(message weboscket.UserSendDirectMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserSendDirectMessage(message, ws, app.userService, app.roomService, app.messageService)
	})

	// test handler
	app.wsServer.SubscribeOnUserSendRoomMessage(func(message weboscket.UserSendRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserSendRoomMessage(message, ws, app.userService, app.roomService)
	})

	log.Println("Server started")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
