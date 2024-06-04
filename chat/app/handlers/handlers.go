package handlers

import (
	"chat/app/components/auth"
	"chat/app/components/message"
	"chat/app/components/room"
	"chat/app/components/user"
	"chat/app/weboscket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func HandleUserConnect(message weboscket.UserConnectMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService) {
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

	log.Println("User connected")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "You are connected")
}

func HandleUserAuth(message weboscket.UserAuthMessage, ws weboscket.WebsocketSender, userService *user.UserService, authService *auth.AuthService, roomService *room.RoomService) {
	fmt.Printf("Handler HandleUserAuth received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserById(ctx, message.Payload.UserID)

	if user == nil {
		log.Println("Error authenticating user: user not found")
		return
	}

	isAuthenticated := authService.AuthenticateUser(user, message.Payload.AccessToken)
	if !isAuthenticated {
		log.Println("User is not authenticated")
		ws.SendMessageToUser(user.Id, "Authentication error")
		return
	}

	log.Println("User is successfully authenticated")

	time.Sleep(2 * time.Second)

	ws.SendMessageToUser(user.Id, "Authentication success!!!")

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

func HandleUserJoinToRoom(message weboscket.UserJoinToRoomMessage, ws weboscket.WebsocketSender, userService *user.UserService, roomService *room.RoomService) {
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

	log.Println("User joined to room")

	time.Sleep(2 * time.Second)

	userName := fmt.Sprintf("Dear user %s", user.Name)
	userMsg := fmt.Sprintf("%s you are joined to room %s", userName, room.Name)

	ws.SendMessageToUser(user.Id, userMsg)

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
		log.Println("Error getting room messages: user not found")
		return
	}

	room := roomService.FindRoomById(message.Payload.RoomID)

	if room == nil {
		log.Println("Error getting room messages: room not found")
		return
	}

	if !room.IsHasUser(user.Id) {
		log.Println("Error getting room messages: user is not in room")
		return
	}

	messages := messageService.FindRoomMessages(ctx, room.Id)

	time.Sleep(2 * time.Second)

	jsonMessages, _ := json.Marshal(messages)

	log.Println("User got room messages:", string(jsonMessages))

	ws.SendMessageToUser(user.Id, string(jsonMessages))
}
