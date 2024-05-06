package message

type MessageService struct {
	messageRepository *MessageRepository
}

func NewMessageService(messageRepository *MessageRepository) *MessageService {
	return &MessageService{messageRepository: messageRepository}
}

func (messageService *MessageService) CreateMessage(text string, creatorID int, receiverID int, roomID int) *Message {
	return NewMessage(1, text, creatorID, receiverID, roomID)
}

func (messageService *MessageService) FindRoomMessages(roomID int) []*Message {
	return messageService.messageRepository.FindMessagesByRoomID(roomID)
}

func (messageService *MessageService) SaveMessage(message *Message) {
	messageService.messageRepository.SaveMessage(message)
}
