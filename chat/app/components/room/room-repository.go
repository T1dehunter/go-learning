package room

type RoomRepository struct {
}

func NewRoomRepository() *RoomRepository {
	return &RoomRepository{}
}

func (roomRepository *RoomRepository) FindRoomById(id int) *Room {
	rooms := make(map[int]*Room)

	rooms[1] = NewRoom(1, "Room 1", Group, []int{1, 2, 3})

	room, ok := rooms[id]
	if !ok {
		return nil
	}

	return room
}
