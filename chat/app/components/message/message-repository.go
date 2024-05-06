package message

import "go.mongodb.org/mongo-driver/mongo"

type MessageRepository struct {
	client *mongo.Client
}

func NewMessageRepository(client *mongo.Client) *MessageRepository {
	return &MessageRepository{client: client}
}

func (messageRepository *MessageRepository) FindMessagesByRoomID(roomID int) []*Message {
	messages := make([]*Message, 3)

	messages[0] = NewMessage(1, "Message 1", 1, 2, 1)
	messages[1] = NewMessage(1, "Message 2", 1, 2, 1)
	messages[2] = NewMessage(1, "Message 3", 1, 2, 1)

	return messages
}

func (messageRepository *MessageRepository) SaveMessage(message *Message) {
	// Save message
}
