package handlers

import (
	"chat/app/components/auth"
	"chat/app/components/room"
	"chat/app/components/user"
	"chat/app/weboscket"
	"fmt"
	"log"
	"time"
)

func HandleUserConnect(message weboscket.UserConnectMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService) {
	fmt.Printf("Handler HandleUserConnect received message -> %+v\n", message)

	user := userService.FindUserById(message.Payload.UserID)

	if user == nil {
		log.Println("Error connecting user: user not found")
		return
	}

	isAuthenticated := authService.AuthenticateUser(user, message.Payload.AccessToken)

	if !isAuthenticated {
		log.Println("Error connecting user: user is not authenticated")
		ws.SendMessageToUser(user.Id, "You are not authenticated")
		return
	}

	log.Println("User connected")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "You are connected")
}

func HandleUserAuth(message weboscket.UserAuthMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserAuth received message -> %+v\n", message)

	user := userService.FindUserById(message.Payload.UserID)

	if user == nil {
		log.Println("Error authenticating user: user not found")
		return
	}

	isAuthenticated := authService.AuthenticateUser(user, message.Payload.AccessToken)
	if !isAuthenticated {
		log.Println("User is not authenticated")
		ws.SendMessageToUser(user.Id, "Authentication error")
		return
	}

	log.Println("User is successfully authenticated")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "Authentication success!!!")

}

func HandleUserJoinToRoom(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserJoinToRoom received message -> %+v\n", message)

	user := userService.FindUserById(message.Payload.UserID)

	if user == nil {
		log.Println("Error joining to room: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)

	if room == nil {
		log.Println("Error joining to room: room not found")
		return
	}

	roomService.JoinUser(user.Id, room.Id)

	log.Println("User joined to room")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "You joined to room")

}

func HandleUserLeaveRoom(message weboscket.UserLeaveRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserLeaveRoom received message -> %+v\n", message)

	user := userService.FindUserById(message.Payload.UserID)

	if user == nil {
		log.Println("Error leaving room: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)

	if room == nil {
		log.Println("Error leaving room: room not found")
		return
	}

	roomService.LeaveUser(user.Id, room.Id)

	log.Println("User left room")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "You left room")
}
