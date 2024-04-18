package weboscket

type UserConnectMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}
}

type UserAuthMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}

type UserCreateDirectRoomMessage struct {
	Name    string `json:"name"`
	Payload struct {
		CreatorID int `json:"creatorID"`
		InviteeID int `json:"inviteeID"`
	} `json:"payload"`
}

type UserJoinToRoomMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserLeaveRoomMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserSendDirectMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID     int    `json:"userID"`
		ReceiverID int    `json:"receiverID"`
		RoomID     int    `json:"roomID"`
		Message    string `json:"message"`
	} `json:"payload"`
}

type UserSendRoomMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID  int    `json:"userID"`
		RoomID  int    `json:"roomID"`
		Message string `json:"message"`
	} `json:"payload"`
}

type WebsocketSender interface {
	SendMessageToUser(connectionID int, message string)
	SendMessageToRoom(connectionID int, message string)
	AddUserToNamespace(namespace string)
	SendMessageToNamespace(namespace string, message string)
}
