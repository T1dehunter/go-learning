package room

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomRepository struct {
	client *mongo.Client
}

func NewRoomRepository(client *mongo.Client) *RoomRepository {
	return &RoomRepository{client: client}
}

func (roomRepository *RoomRepository) FindRoomById(id int) *Room {
	rooms := make(map[int]*Room)

	rooms[1] = NewRoom(1, "Room 1", "group", []int{1, 2, 3})

	room, ok := rooms[id]
	if !ok {
		return nil
	}

	return room
}

func (roomRepository *RoomRepository) DeleteAllRooms(ctx context.Context) {
	roomRepository.client.Database("chat").Collection("rooms").Drop(ctx)
}

func (roomRepository *RoomRepository) AddRoom(ctx context.Context, room *Room) {
	collection := roomRepository.client.Database("chat").Collection("rooms")
	data := map[string]interface{}{"id": room.Id, "name": room.Name, "type": room.Type, "userIds": room.UserIds}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
}
