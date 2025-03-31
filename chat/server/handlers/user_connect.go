package handlers

import (
	"chat/server/components/auth"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/handlers/messages"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"log"
	"time"
)

func HandleUserConnect(
	message weboscket.UserConnectMsg,
	userService *user.UserService,
	authService *auth.AuthService,
	roomService *room.RoomService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)
	if user == nil {
		log.Println("Error connecting user: user not found")
		return
	}

	isAuthenticated := authService.AuthenticateUser(user, message.Payload.AccessToken)
	if !isAuthenticated {
		log.Println("Error connecting user: user is not authenticated")
		response.SendMessageToUser(user.Id, "You are not authenticated")
		return
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	userRooms := roomService.FindUserRooms(ctx, user.Id)

	roomsResponse := []messages.UserRoom{}
	for _, room := range userRooms {
		users := userService.FindAllUsersByIds(ctx, room.UserIds)
		var roomUsers []messages.UserData
		for _, user := range users {
			roomUsers = append(roomUsers, messages.UserData{
				ID:   user.Id,
				Name: user.Name,
			})
		}
		roomsResponse = append(roomsResponse, messages.UserRoom{
			ID:    room.Id,
			Name:  room.Name,
			Users: roomUsers,
		})
	}

	msg := messages.UserConnectedMsg{
		Type: "user_connected",
		Payload: struct {
			Success bool                `json:"success"`
			Rooms   []messages.UserRoom `json:"rooms"`
		}{
			Success: true,
			Rooms:   roomsResponse,
		},
	}

	msgJson, _ := json.Marshal(msg)

	response.SendMessageToUser(user.Id, string(msgJson))
}
