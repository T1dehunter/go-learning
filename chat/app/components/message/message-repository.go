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

func (messageRepository *MessageRepository) FindMessagesByRoomID(ctx context.Context, roomID int) []*Message {
	collection := messageRepository.client.Database("chat").Collection("messages")
	cursor, err := collection.Find(ctx, map[string]interface{}{"roomID": roomID})
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	messages := []*Message{}
	for cursor.Next(ctx) {
		var msg Message
		err := cursor.Decode(&msg)
		if err != nil {
			panic(err)
		}
		messages = append(messages, NewMessage(msg.Id, msg.Text, msg.CreatorID, msg.ReceiverID, msg.RoomID))
	}
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
