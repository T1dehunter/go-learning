package handlers

import (
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"log"
	"time"
)

func HandleGetRoomMessages(
	message weboscket.UserGetListRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	messageService *message.MessageService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)
	if user == nil {
		log.Println("Error getting room events: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)
	if room == nil {
		log.Println("Error getting room events: room not found")
		return
	}
	if !room.IsHasUser(user.Id) {
		log.Println("Error getting room events: user is not in room")
		return
	}

	messages := messageService.FindRoomMessages(ctx, room.Id)

	messagesJson, _ := json.Marshal(messages)

	response.SendMessageToUser(user.Id, string(messagesJson))
}
