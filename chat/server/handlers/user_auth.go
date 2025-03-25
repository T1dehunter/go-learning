package handlers

import (
	"chat/server/components/user"
	"chat/server/handlers/messages"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

func HandleUserAuth(
	message weboscket.UserAuthMsg,
	userService *user.UserService,
	response *weboscket.Response,
) {
	fmt.Printf("Handler HandleUserAuth received message -> %+v\n", message)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserByNameAndPassword(ctx, message.Payload.UserName, message.Payload.Password)
	if user == nil {
		log.Println("Error authenticating user: user not found")
		res := messages.UserAuthError{
			Type: "user_not_authenticated",
		}
		resMsg, err := json.Marshal(res)
		if err != nil {
			fmt.Println("Error converting messages to JSON:", err)
			return
		}
		response.SendMessageToUser(user.Id, string(resMsg))
		return
	}

	log.Println("User is successfully authenticated")

	time.Sleep(2 * time.Second)

	// for testing purposes only, user's password is sent back to the user as access token
	//responseMsg := fmt.Sprintf("{\"type\": \"user_authenticated\", \"payload\": {\"userID\": %d, \"userName\": \"%s\", \"accessToken\": \"%s\"}}", user.Id, user.Name, user.Password)

	res := messages.UserAuthenticatedMsg{
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
		fmt.Println("Error converting messages to JSON:", err)
		return
	}

	response.SendMessage(string(resMsg))
}
