package weboscket

type UserAuthMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	} `json:"payload"`
}

type UserConnectMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}
}

type UserCreateDirectRoomMessage struct {
	Type    string `json:"type"`
	Payload struct {
		CreatorID int `json:"creatorID"`
		InviteeID int `json:"inviteeID"`
	} `json:"payload"`
}

type UserJoinToRoomMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserLeaveRoomMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserSendDirectMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserID     int    `json:"userID"`
		ReceiverID int    `json:"receiverID"`
		RoomID     int    `json:"roomID"`
		Message    string `json:"message"`
	} `json:"payload"`
}

type UserSendRoomMessage struct {
	Type    string `json:"type"`
	Payload struct {
		UserID  int    `json:"userID"`
		RoomID  int    `json:"roomID"`
		Message string `json:"message"`
	} `json:"payload"`
}

type UserGetRoomMessages struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type WebsocketSender interface {
	RegisterConnection(userID int)
	SendMessageToUser(connectionID int, message string)
	SendMessageToRoom(connectionID int, message string)
	AddUserToNamespace(namespace string)
	SendMessageToNamespace(namespace string, message string)
}
