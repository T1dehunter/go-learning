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

func HandleUserSendDirectMessage(
	message weboscket.UserSendDirectMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	messageService *message.MessageService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)
	if user == nil {
		log.Println("Error sending direct message: user not found")
		return
	}

	receiver := userService.FindUserById(ctx, message.Payload.ReceiverID)
	if receiver == nil {
		log.Println("Error sending direct message: receiver not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)
	if room == nil {
		log.Println("Error sending direct message: room not found")
		return
	}

	if !room.IsHasUser(user.Id) {
		log.Println("Error sending direct message: user is not in room")
		return
	}

	if !room.IsHasUser(receiver.Id) {
		log.Println("Error sending direct message: receiver is not in room")
		return
	}

	now := time.Now()
	newMsg := messageService.CreateMessage(
		message.Payload.Message,
		user.Id,
		receiver.Id,
		room.Id,
		now.String(),
	)

	messageService.SaveMessage(newMsg)

	newMsgJson, _ := json.Marshal(newMsg)

	response.SendMessageToUser(receiver.Id, string(newMsgJson))
}
