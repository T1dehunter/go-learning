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
	room.joinUser(userId)
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
