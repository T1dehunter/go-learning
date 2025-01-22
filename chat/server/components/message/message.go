package message

type Message struct {
	Id         int
	Text       string
	CreatorID  int
	ReceiverID int
	RoomID     int
	CreatedAt  string
}

func NewMessage(id int, text string, creatorID int, receiverID int, roomID int, createdAt string) *Message {
	return &Message{Id: id, Text: text, CreatorID: creatorID, ReceiverID: receiverID, RoomID: roomID, CreatedAt: createdAt}
}
