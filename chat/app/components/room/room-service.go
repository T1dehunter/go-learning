package room

import "fmt"

type RoomService struct {
	roomRepository *RoomRepository
}

func NewRoomService() *RoomService {
	return &RoomService{}
}

func (romService *RoomService) FindRoomById(id int) *Room {
	return romService.roomRepository.FindRoomById(id)
}

func (roomService *RoomService) JoinUser(userId int, roomId int) bool {
	room := roomService.FindRoomById(roomId)
	if room == nil {
		fmt.Println("Error joining user to room: room not found")
		return false
	}
	room.JoinUser(userId)
	return true
}

func (roomService *RoomService) LeaveUser(userId int, roomId int) bool {
	room := roomService.FindRoomById(roomId)
	if room == nil {
		fmt.Println("Error leaving user from room: room not found")
		return false
	}
	room.leaveUser(userId)
	return true
}

func (roomService *RoomService) CreateDirectRoom(name string) *Room {
	return NewRoom(1, name, Direct, []int{})
}

func (roomService *RoomService) SaveRoom(room Room) bool {
	return true
}
