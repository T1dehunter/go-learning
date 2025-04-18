package handlers

import (
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/handlers/messages"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"log"
	"time"
)

func HandleUserJoinRoom(
	message weboscket.UserJoinToRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	messageService *message.MessageService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)
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

	roomMessages := messageService.FindRoomMessages(ctx, room.Id)

	roomUsers := userService.FindAllUsersByIds(ctx, room.UserIds)
	var roomUsersData []messages.UserData
	for _, user := range roomUsers {
		roomUsersData = append(roomUsersData, messages.UserData{
			ID:   user.Id,
			Name: user.Name,
		})
	}

	var roomMessagesData []messages.Message
	for _, message := range roomMessages {
		var creatorUser messages.UserData
		for _, user := range roomUsersData {
			if user.ID == message.CreatorID {
				creatorUser = messages.UserData{
					ID:   user.ID,
					Name: user.Name,
				}
			}
		}

		var receiverUser messages.UserData
		for _, user := range roomUsersData {
			if user.ID == message.ReceiverID {
				receiverUser = messages.UserData{
					ID:   user.ID,
					Name: user.Name,
				}
			}
		}

		roomMessagesData = append(roomMessagesData, messages.Message{
			ID:           message.Id,
			CreatorID:    message.CreatorID,
			CreatorName:  creatorUser.Name,
			ReceiverID:   message.ReceiverID,
			ReceiverName: receiverUser.Name,
			RoomID:       message.RoomID,
			Text:         message.Text,
			CreatedAt:    message.CreatedAt,
		})
	}

	msg := messages.UserJoinedToRoomMsg{
		Type: "user_joined_to_room",
		Payload: struct {
			Success  bool                `json:"success"`
			RoomID   int                 `json:"roomID"`
			RoomName string              `json:"roomName"`
			Users    []messages.UserData `json:"users"`
			Messages []messages.Message  `json:"messages"`
		}{
			Success:  true,
			RoomID:   room.Id,
			RoomName: room.Name,
			Users:    roomUsersData,
			Messages: roomMessagesData,
		},
	}

	msgJson, _ := json.Marshal(msg)

	response.SendMessageToUser(user.Id, string(msgJson))
}
