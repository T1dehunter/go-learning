package types

type User struct {
	ID   int
	Name string
}

type Room struct {
	ID    int
	Name  string
	Type  string
	Users []User
}

type Message struct {
	ID          int
	RoomID      int
	CreatorID   int
	CreatorName string
	ReceiverID  int
	Text        string
	CreatedAt   string
}
