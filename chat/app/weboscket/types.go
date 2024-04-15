package weboscket

type UserAuthMessage struct {
	Name    string `json:"name"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
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

type WebsocketSender interface {
	SendMessageToUser(connectionID int, message string)
}
