package handlers

import (
	"chat/server/components/user"
	"chat/server/handlers/messages"
	"chat/server/weboscket"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func HandleUserAuth(
	message weboscket.UserAuthMsg,
	userService *user.UserService,
	response *weboscket.Response,
) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	user := userService.FindUserByNameAndPassword(ctx, message.Payload.UserName, message.Payload.Password)
	if user == nil {
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

	msg := messages.UserAuthenticatedMsg{
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

	msgJson, _ := json.Marshal(msg)

	response.SendMessage(string(msgJson))
}
