package room

type RoomService struct {
	getRoomById func(id int) string
	joinUser    func(roomId int, userId int) string
}

func NewRoomService() *RoomService {
	return &RoomService{
		getRoomById: func(id int) string {
			return "Room"
		},
		joinUser: func(roomId int, userId int) string {
			return "User joined room"
		},
	}
}
