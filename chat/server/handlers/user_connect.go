package handlers

import (
	"chat/server/components/auth"
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

func HandleUserConnect(
	message weboscket.UserConnectMsg,
	userService *user.UserService,
	authService *auth.AuthService,
	roomService *room.RoomService,
	response *weboscket.Response,
) {
	fmt.Printf("Handler HandleUserConnect received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	fmt.Println("HandleUserConnect message.Payload.UserID", message.Payload)

	user := userService.FindUserById(ctx, message.Payload.UserID)
	if user == nil {
		log.Println("Error connecting user: user not found")
		return
	}
	fmt.Println("user user user: ", user.Password)
	isAuthenticated := authService.AuthenticateUser(user, message.Payload.AccessToken)
	if !isAuthenticated {
		log.Println("Error connecting user: user is not authenticated")
		response.SendMessageToUser(user.Id, "You are not authenticated")
		return
	}

	//ws.RegisterConnection(user.Id)

	//ws.AddUserToNamespace("connected_users")

	if user.RoomID != nil {
		//roomName := fmt.Sprintf("room_%d", *user.RoomID)
		//ws.AddUserToNamespace(roomName)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	userRooms := roomService.FindUserRooms(ctx, user.Id)

	time.Sleep(2 * time.Second)

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
	res := messages.UserConnectedMsg{
		Type: "user_connected",
		Payload: struct {
			Success bool                `json:"success"`
			Rooms   []messages.UserRoom `json:"rooms"`
		}{
			Success: true,
			Rooms:   roomsResponse,
		},
	}
	resMsg, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error converting messages to JSON:", err)
		return
	}
	fmt.Println("resMsg", string(resMsg))
	response.SendMessageToUser(user.Id, string(resMsg))
}
