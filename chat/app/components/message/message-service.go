package message

type MessageService struct {
	roomRepository *MessageRepository
}

func NewMessageService() *MessageService {
	return &MessageService{roomRepository: NewMessageRepository()}
}

func (messageService *MessageService) CreateMessage(text string, creatorID int, receiverID int, roomID int) *Message {
	return NewMessage(1, text, creatorID, receiverID, roomID)
}

func (messageService *MessageService) FindRoomMessages(roomID int) []*Message {
	return messageService.roomRepository.FindMessagesByRoomID(roomID)
}

func (messageService *MessageService) SaveMessage(message *Message) {
	messageService.roomRepository.SaveMessage(message)
}
