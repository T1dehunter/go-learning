package response

type UserAuthenticatedMsg struct {
	Type    string `json:"type"`
	Payload struct {
		UserID      int    `json:"userID"`
		UserName    string `json:"userName"`
		AccessToken string `json:"accessToken"`
	} `json:"payload"`
}

type UserData struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type UserRoom struct {
	ID    int        `json:"id"`
	Name  string     `json:"name"`
	Type  string     `json:"type"`
	Users []UserData `json:"users"`
}

type UserConnectedMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Success bool       `json:"success"`
		Rooms   []UserRoom `json:"rooms"`
	} `json:"payload"`
}

type Message struct {
	ID          int    `json:"id"`
	RoomID      int    `json:"roomID"`
	CreatorID   int    `json:"CreatorID"`
	CreatorName string `json:"CreatorName"`
	ReceiverID  int    `json:"receiverID"`
	Text        string `json:"text"`
	CreatedAt   string `json:"createdAt"`
}
type UserJoinedToRoomMsg struct {
	Type    string `json:"type"`
	Payload struct {
		Success  bool       `json:"success"`
		RoomID   int        `json:"roomID"`
		RoomName string     `json:"roomName"`
		Users    []UserData `json:"users"`
		Messages []Message  `json:"messages"`
	} `json:"payload"`
}
