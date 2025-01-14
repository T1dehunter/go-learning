package handlers

import (
	"chat/server/components/auth"
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/handlers/response"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func HandleUserAuth(message weboscket.UserAuthMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserAuth received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserByNameAndPassword(ctx, message.Payload.UserName, message.Payload.Password)
	if user == nil {
		log.Println("Error authenticating user: user not found")
		responseMsg := fmt.Sprintf("{\"type\": \"user_not_authenticated\", \"payload\": {\"userID\": \"\", \"accessToken\": \"\"}}")
		ws.SendMessageToUser(1, responseMsg)
		return
	}

	log.Println("User is successfully authenticated")

	time.Sleep(2 * time.Second)

	// for testing purposes only, user's password is sent back to the user as access token
	//responseMsg := fmt.Sprintf("{\"type\": \"user_authenticated\", \"payload\": {\"userID\": %d, \"userName\": \"%s\", \"accessToken\": \"%s\"}}", user.Id, user.Name, user.Password)

	res := response.UserAuthenticatedMsg{
		Type: "user_authenticated",
		Payload: struct {
			UserID      int    `json:"userID"`
			UserName    string `json:"userName"`
			AccessToken string `json:"accessToken"`
		}{
			UserID:      user.Id,
			UserName:    user.Name,
			AccessToken: user.Password,
		},
	}

	resMsg, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error converting response to JSON:", err)
		return
	}
	ws.SendMessageToUser(user.Id, string(resMsg))

}

func HandleUserConnect(message weboscket.UserConnectMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService, roomService *room.RoomService) {
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
		ws.SendMessageToUser(user.Id, "You are not authenticated")
		return
	}

	ws.RegisterConnection(user.Id)

	ws.AddUserToNamespace("connected_users")

	if user.RoomID != nil {
		roomName := fmt.Sprintf("room_%d", *user.RoomID)
		ws.AddUserToNamespace(roomName)
	}

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	userRooms := roomService.FindUserRooms(ctx, user.Id)

	time.Sleep(2 * time.Second)

	roomsResponse := []response.UserRoom{}
	for _, room := range userRooms {
		users := userService.FindAllUsersByIds(ctx, room.UserIds)
		var roomUsers []response.UserData
		for _, user := range users {
			roomUsers = append(roomUsers, response.UserData{
				ID:   user.Id,
				Name: user.Name,
			})
		}
		roomsResponse = append(roomsResponse, response.UserRoom{
			ID:    room.Id,
			Name:  room.Name,
			Users: roomUsers,
		})
	}
	res := response.UserConnectedMsg{
		Type: "user_connected",
		Payload: struct {
			Success bool                `json:"success"`
			Rooms   []response.UserRoom `json:"rooms"`
		}{
			Success: true,
			Rooms:   roomsResponse,
		},
	}
	resMsg, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error converting response to JSON:", err)
		return
	}
	fmt.Println("resMsg", string(resMsg))
	ws.SendMessageToUser(user.Id, string(resMsg))
}

func HandleUserCreateDirectRoom(message weboscket.UserCreateDirectRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
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

	ws.SendMessageToUser(roomCreator.Id, "Direct room created")

	log.Println("Direct room created")

}

func HandleUserJoinToRoom(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService, messageService *message.MessageService) {
	fmt.Printf("Handler HandleUserJoinToRoom received message -> %+v\n", message)

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

	roomName := fmt.Sprintf("room_%d", room.Id)
	ws.AddUserToNamespace(roomName)

	roomMessages := messageService.FindRoomMessages(ctx, room.Id)

	roomUsers := userService.FindAllUsersByIds(ctx, room.UserIds)
	var roomUsersData []response.UserData
	for _, user := range roomUsers {
		roomUsersData = append(roomUsersData, response.UserData{
			ID:   user.Id,
			Name: user.Name,
		})
	}

	var roomMessagesData []response.Message
	for _, message := range roomMessages {
		roomMessagesData = append(roomMessagesData, response.Message{
			ID:         message.Id,
			CreatorID:  message.CreatorID,
			ReceiverID: message.ReceiverID,
			RoomID:     message.RoomID,
			Text:       message.Text,
			CreatedAt:  message.CreatedAt,
		})
	}

	log.Println("User joined to room")

	time.Sleep(2 * time.Second)

	userName := fmt.Sprintf("Dear user %s", user.Name)
	userMsg := fmt.Sprintf("%s you are joined to room %s", userName, room.Name)

	fmt.Println("userMsg", userMsg)

	res := response.UserJoinedToRoomMsg{
		Type: "user_joined_to_room",
		Payload: struct {
			Success  bool                `json:"success"`
			RoomID   int                 `json:"roomID"`
			RoomName string              `json:"roomName"`
			Users    []response.UserData `json:"users"`
			Messages []response.Message  `json:"messages"`
		}{
			Success:  true,
			RoomID:   room.Id,
			RoomName: room.Name,
			Users:    roomUsersData,
			Messages: roomMessagesData,
		},
	}
	resMsg, err := json.Marshal(res)
	if err != nil {
		fmt.Println("Error converting response to JSON:", err)
		return
	}
	ws.SendMessageToUser(user.Id, string(resMsg))
}

func HandleUserLeaveRoom(message weboscket.UserLeaveRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
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

	ws.SendMessageToUser(user.Id, "You left room")
}

func HandleUserSendDirectMessage(message weboscket.UserSendDirectMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService, messageService *message.MessageService) {
	fmt.Printf("Handler HandleUserSendDirectMessage received message -> %+v\n", message)

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

	newMessage := messageService.CreateMessage(message.Payload.Message, user.Id, receiver.Id, room.Id)

	messageService.SaveMessage(newMessage)

	log.Println("User sent direct message")

	time.Sleep(2 * time.Second)

	jsonMessage, _ := json.Marshal(newMessage)

	ws.SendMessageToUser(receiver.Id, string(jsonMessage))
}

func HandleUserSendRoomMessage(message weboscket.UserSendRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserSendRoomMessage received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)

	if user == nil {
		log.Println("Error sending room message: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)

	if room == nil {
		log.Println("Error sending room message: room not found")
		return
	}

	log.Println("User sent message to room")

	time.Sleep(2 * time.Second)

	roomName := fmt.Sprintf("room_%d", room.Id)

	ws.SendMessageToNamespace(roomName, message.Payload.Message)
}

func HandleGetRoomMessages(message weboscket.UserGetRoomMessages, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService, messageService *message.MessageService) {
	fmt.Printf("Handler HandleGetRoomMessages received message -> %+v\n", message)

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

	time.Sleep(2 * time.Second)

	jsonMessages, _ := json.Marshal(messages)

	log.Println("User got room events:", string(jsonMessages))

	ws.SendMessageToUser(user.Id, string(jsonMessages))
}
