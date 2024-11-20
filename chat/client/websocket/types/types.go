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

type Room struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}
type UserConnectMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool   `json:"success"`
		Rooms   []Room `json:"rooms"`
	}
}

type UserJoinToRoomMessageWs struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		RoomID      int    `json:"roomID"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}

type UserJoinToRoomMessageResponseWs struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool   `json:"success"`
		RoomID   int    `json:"roomID"`
		RoomName string `json:"roomName"`
		Msg      string `json:"msg"`
	}
}
