package response

type UserAuthenticatedMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}

type UserRoom struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type UserConnectedMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool       `json:"success"`
		Rooms   []UserRoom `json:"rooms"`
	} `json:"payload"`
}

type UserJoinedToRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool   `json:"success"`
		RoomID   int    `json:"roomID"`
		RoomName string `json:"roomName"`
		Msg      string `json:"msg"`
	} `json:"payload"`
}
