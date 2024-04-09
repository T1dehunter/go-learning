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
	UserService *user.UserService
	AuthService *auth.AuthService
	RoomService *room.RoomService
	WebSocket   *weboscket.WebSocket
}

func NewApp() *App {
	return &App{
		UserService: user.NewUserService(),
		AuthService: auth.NewAuthService(),
		RoomService: room.NewRoomService(),
		WebSocket:   weboscket.NewWebSocket(),
	}
}

func (app *App) Start() {
	http.HandleFunc("/chat", app.WebSocket.Listen)

	app.WebSocket.SubscribeOnUserAuth(func(message weboscket.UserAuthMessage) {
		handlers.HandleUserAuth(message, app.UserService, app.AuthService, app.RoomService)
	})

	log.Println("Server started")

	log.Fatal(http.ListenAndServe(":3000", nil))
}
