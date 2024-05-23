package message

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

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

func (messageRepository *MessageRepository) DeleteAllMessages(ctx context.Context) {
	messageRepository.client.Database("chat").Collection("messages").Drop(ctx)
}

func (messageRepository *MessageRepository) AddMessage(ctx context.Context, message *Message) {
	collection := messageRepository.client.Database("chat").Collection("messages")
	data := map[string]interface{}{
		"id":         message.Id,
		"creatorID":  message.CreatorID,
		"receiverID": message.ReceiverID,
		"roomID":     message.RoomID,
		"text":       message.Text,
	}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
}
