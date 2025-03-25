package handlers

import (
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/weboscket"
	"context"
	"fmt"
	"log"
	"time"
)

func HandleUserLeaveRoom(
	message weboscket.UserLeaveRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	response *weboscket.Response,
) {
	fmt.Printf("Handler HandleUserLeaveRoom received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)

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

	response.SendMessageToUser(user.Id, "You left room")
}
