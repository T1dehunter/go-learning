package handlers

import (
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/weboscket"
	"context"
	"log"
	"time"
)

func HandleUserLeaveRoom(
	message weboscket.UserLeaveRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)
	if user == nil {
		log.Println("Error on leaving room: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)
	if room == nil {
		log.Println("Error on leaving room: room not found")
		return
	}

	roomService.LeaveUser(user.Id, room.Id)

	response.SendMessageToUser(user.Id, "You left room")
}
