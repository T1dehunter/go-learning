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
	collection := roomRepository.client.Database("chat").Collection("rooms")
	data := map[string]interface{}{"id": id}
	var room Room
	err := collection.FindOne(context.Background(), data).Decode(&room)
	if err != nil {
		return nil
	}
	return &room
}

func (roomRepository *RoomRepository) FindRoomsByUserId(ctx context.Context, userId int) []*Room {
	rooms := make([]*Room, 0)

	collection := roomRepository.client.Database("chat").Collection("rooms")
	// map with inner map
	data := map[string]interface{}{"userIds": map[string]interface{}{"$in": []int{userId}}}
	cursor, err := collection.Find(ctx, data)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	rooms = make([]*Room, 0)
	for cursor.Next(ctx) {
		var room Room
		err := cursor.Decode(&room)
		if err != nil {
			panic(err)
		}
		rooms = append(rooms, &room)
	}

	return rooms
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
