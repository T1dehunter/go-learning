package handlers

import (
	"chat/app/components/auth"
	"chat/app/components/room"
	"chat/app/components/user"
	"chat/app/weboscket"
	"log"
)

func HandleUserAuth(message weboscket.UserAuthMessage, userService *user.UserService, authService *auth.AuthService, roomService *room.RoomService) {
	authService.AuthenticateUser(message.Payload.UserID, message.Payload.AccessToken)
	log.Println("HANDLER HandleUserAuth received message:", message)
}
