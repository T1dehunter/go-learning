package weboscket

type UserAuthMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	} `json:"payload"`
}

type UserConnectMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}
}

type UserCreateRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		CreatorID int `json:"creatorID"`
		InviteeID int `json:"inviteeID"`
	} `json:"payload"`
}

type UserJoinToRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserSendRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID  int    `json:"userID"`
		RoomID  int    `json:"roomID"`
		Message string `json:"message"`
	} `json:"payload"`
}

type UserLeaveRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type UserSendDirectMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID     int    `json:"userID"`
		ReceiverID int    `json:"receiverID"`
		RoomID     int    `json:"roomID"`
		Message    string `json:"message"`
	} `json:"payload"`
}

type UserGetListRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	} `json:"payload"`
}

type ClientLogMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Text string `json:"text"`
	} `json:"payload"`
}

type LogMsg struct {
	Type      string `json:"type"`
	Title     string `json:"title"`
	CreatedAt string `json:"createdAt"`
}

type WebsocketSender interface {
	SendMessageToUser(connectionID int, message string)
	SendLog(msg LogMsg)
}
