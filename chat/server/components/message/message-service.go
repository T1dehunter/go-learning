package message

import "context"

type MessageService struct {
	messageRepository *MessageRepository
}

func NewMessageService(messageRepository *MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (messageService *MessageService) CreateMessage(text string, creatorID int, receiverID int, roomID int, createdAt string) *Message {
	return NewMessage(1, text, creatorID, receiverID, roomID, createdAt)
}

func (messageService *MessageService) FindRoomMessages(ctx context.Context, roomID int) []*Message {
	return messageService.messageRepository.FindMessagesByRoomID(ctx, roomID)
}

func (messageService *MessageService) SaveMessage(message *Message) {
	messageService.messageRepository.AddMessage(nil, message)
}
