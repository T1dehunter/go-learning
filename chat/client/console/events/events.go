package events

type UserChatConfirmed struct {
	IsConfirmed bool
}

type UserAuthRequest struct {
	Username string
	Password string
}

type UserAuthFailedRes struct {
}

type UserAuthSuccessRes struct {
	Payload struct {
		UserID      int
		UserName    string
		AccessToken string
	}
}

type Room struct {
	ID   int
	Name string
}

type UserConnectRes struct {
	Payload struct {
		UserID      int
		UserName    string
		AccessToken string
		Rooms       []Room
	}
}

type UserJoinRoom struct {
	RoomID int
}

type UserSendRoomMessage struct {
	RoomID  int
	Message string
}

type UserRoomExit struct{}

type UserChatExit struct{}
