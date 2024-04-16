package app

import (
	"chat/app/components/auth"
	"chat/app/components/room"
	"chat/app/components/user"
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
	userService *user.UserService
	authService *auth.AuthService
	roomService *room.RoomService
	wsServer    *weboscket.WebSocketServer
}

func NewApp() *App {
	return &App{
		userService: user.NewUserService(),
		authService: auth.NewAuthService(),
		roomService: room.NewRoomService(),
		wsServer:    weboscket.NewWebSocketServer(),
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

	app.wsServer.SubscribeOnUserJoinToRoom(func(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserJoinToRoom(message, ws, app.userService, app.roomService)
	})

	app.wsServer.SubscribeOnUserLeaveRoom(func(message weboscket.UserLeaveRoomMessage, ws weboscket.WebsocketSender) {
		handlers.HandleUserLeaveRoom(message, ws, app.userService, app.roomService)
	})

	log.Println("Server started")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
