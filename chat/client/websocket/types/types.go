package types

type UserAuthMessageWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
}

type UserAuthMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	}
}

type UserConnectMessageWs struct {
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
type UserConnectMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool   `json:"success"`
		Rooms   []Room `json:"rooms"`
	} `json:"payload"`
}

type UserJoinToRoomMessageWs struct {
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
type UserJoinToRoomMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool      `json:"success"`
		RoomID   int       `json:"roomID"`
		RoomName string    `json:"roomName"`
		Users    []User    `json:"users"`
		Messages []Message `json:"messages"`
	}
}

type UserSendRoomMessageWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		RoomID      int    `json:"roomID"`
		Message     string `json:"message"`
		AccessToken string `json:"accessToken"`
	}
}

type UserSendRoomMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool      `json:"success"`
		RoomID   int       `json:"roomID"`
		RoomName string    `json:"roomName"`
		Users    []User    `json:"users"`
		Messages []Message `json:"messages"`
	}
}
