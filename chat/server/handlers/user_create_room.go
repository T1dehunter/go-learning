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

func HandleUserCreateRoom(
	message weboscket.UserCreateRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	response *weboscket.Response,
) {
	fmt.Printf("Handler HandleUserCreateDirectRoom received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	roomCreator := userService.FindUserById(ctx, message.Payload.CreatorID)
	if roomCreator == nil {
		log.Println("Error creating direct room: creator not found")
		return
	}

	roomInvitee := userService.FindUserById(ctx, message.Payload.InviteeID)
	if roomInvitee == nil {
		log.Println("Error creating direct room: invitee not found")
		return

	}

	room := roomService.CreateDirectRoom("Direct room")
	room.JoinUser(roomCreator.Id)
	room.JoinUser(roomInvitee.Id)

	roomService.SaveRoom(*room)

	response.SendMessageToUser(roomCreator.Id, "Direct room created")

	log.Println("Direct room created")
}
