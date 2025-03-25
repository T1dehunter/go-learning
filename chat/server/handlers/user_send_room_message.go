package handlers

import (
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/handlers/messages"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func HandleUserSendRoomMessage(
	message weboscket.UserSendRoomMsg,
	userService *user.UserService,
	roomService *room.RoomService,
	messageService *message.MessageService,
	response *weboscket.RoomResponse,
) {

	logMsg := fmt.Sprintf("[Handler] HandleUserSendRoomMessage received message -> %+v\n", message)
	response.LogText(logMsg)

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

	currentTime := time.Now().UTC()
	createdAt := currentTime.Format("2006-01-02T15:04:05.000Z")
	newMsg := messageService.CreateMessage(message.Payload.Message, user.Id, user.Id, room.Id, createdAt)
	messageService.SaveMessage(newMsg)

	roomMessages := messageService.FindRoomMessages(ctx, room.Id)
	fmt.Printf("MESSAGESSSS %+v\n", roomMessages[0])

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

	log.Println("User joined to room")

	time.Sleep(2 * time.Second)

	userName := fmt.Sprintf("Dear user %s", user.Name)
	userMsg := fmt.Sprintf("%s you are joined to room %s", userName, room.Name)

	fmt.Println("userMsg", userMsg)

	res := messages.UserSendRoomMsg{
		Type: "user_send_room_message",
		Payload: struct {
			Success  bool                `json:"success"`
			RoomID   int                 `json:"roomID"`
			RoomName string              `json:"roomName"`
			UserID   int                 `json:"userID"`
			Users    []messages.UserData `json:"users"`
			Messages []messages.Message  `json:"messages"`
		}{
			Success:  true,
			RoomID:   room.Id,
			RoomName: room.Name,
			UserID:   user.Id,
			Users:    roomUsersData,
			Messages: roomMessagesData,
		},
	}
	resMsg, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error converting messages to JSON:", err)
		return
	}

	response.SendToAll(string(resMsg))
}
