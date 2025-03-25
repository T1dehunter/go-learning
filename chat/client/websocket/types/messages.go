package types

type UserAuthMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
}

type UserAuthMsgResponse struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}
}

type UserConnectMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		AccessToken string `json:"accessToken"`
	}
}

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Room struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Users []User `json:"users"`
}
type UserConnectMsgResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool   `json:"success"`
		Rooms   []Room `json:"rooms"`
	} `json:"payload"`
}

type UserJoinToRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		RoomID      int    `json:"roomID"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}

type Message struct {
	ID           int    `json:"id"`
	RoomID       int    `json:"roomID"`
	CreatorID    int    `json:"creatorID"`
	CreatorName  string `json:"creatorName"`
	ReceiverID   int    `json:"receiverID"`
	ReceiverName string `json:"receiverName"`
	Text         string `json:"text"`
	CreatedAt    string `json:"createdAt"`
}
type UserJoinToRoomMsgResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool      `json:"success"`
		RoomID   int       `json:"roomID"`
		RoomName string    `json:"roomName"`
		Users    []User    `json:"users"`
		Messages []Message `json:"messages"`
	}
}

type UserSendRoomMsgPayload struct {
	UserID      int    `json:"userID"`
	RoomID      int    `json:"roomID"`
	Message     string `json:"message"`
	AccessToken string `json:"accessToken"`
}
type UserSendRoomMsg struct {
	Type    string                 `json:"type"`
	Payload UserSendRoomMsgPayload `json:"payload"`
}

type UserSendRoomMsgResponse struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool      `json:"success"`
		RoomID   int       `json:"roomID"`
		RoomName string    `json:"roomName"`
		UserID   int       `json:"userID"`
		Users    []User    `json:"users"`
		Messages []Message `json:"messages"`
	}
}

type UserLeaveRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID int `json:"userID"`
		RoomID int `json:"roomID"`
	}
}

type ClientLogMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Text string `json:"text"`
	}
}
