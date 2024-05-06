package room

import "go.mongodb.org/mongo-driver/mongo"

type RoomRepository struct {
	client *mongo.Client
}

func NewRoomRepository(client *mongo.Client) *RoomRepository {
	return &RoomRepository{client: client}
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
