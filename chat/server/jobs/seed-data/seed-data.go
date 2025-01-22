package seed_data

import (
	"chat/server/components/message"
	"chat/server/components/room"
	"chat/server/components/user"
	"chat/server/database"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"os"
	"time"
)

type UserJson struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	IsBanned bool   `json:"isBanned"`
	RoomID   *int   `json:"roomID"`
}

type RoomJson struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	UserIds []int  `json:"userIds"`
}

type MessageJson struct {
	Id         int    `json:"id"`
	CreatorID  int    `json:"creatorID"`
	ReceiverID int    `json:"receiverID"`
	RoomID     int    `json:"roomID"`
	Text       string `json:"text"`
	CreatedAt  string `json:"createdAt"`
}

const pathToUsersJson = "server/jobs/seed-data/data/users.json"
const pathToRoomsJson = "server/jobs/seed-data/data/rooms.json"
const pathToMessagesJson = "server/jobs/seed-data/data/messages.json"

func Seed() {
	fmt.Println("Start seeding data...")

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client := database.Connect()

	seedUsers(ctx, client)
	seedRooms(ctx, client)
	seedMessages(ctx, client)
}

func seedUsers(ctx context.Context, client *mongo.Client) {
	userRepository := user.NewUserRepository(client)
	users := readUsers()
	fmt.Println("Deleting current users...")
	userRepository.DeleteAllUsers(ctx)
	fmt.Println("Adding new users...")
	for _, user := range users {
		userRepository.AddUser(ctx, &user)
	}
}

func readUsers() []user.User {
	users := readJsonFile[UserJson](pathToUsersJson)
	if users == nil {
		return nil
	}
	var userEntities []user.User
	for _, userJson := range users {
		user := user.NewUser(userJson.Id, userJson.Name, userJson.Email, userJson.Password, userJson.IsBanned, userJson.RoomID)
		userEntities = append(userEntities, *user)
	}
	return userEntities
}

func seedRooms(ctx context.Context, client *mongo.Client) {
	roomRepository := room.NewRoomRepository(client)
	rooms := readRooms()
	fmt.Println("Deleting current rooms...")
	roomRepository.DeleteAllRooms(ctx)
	fmt.Println("Adding new rooms...")
	for _, room := range rooms {
		roomRepository.AddRoom(ctx, &room)
	}
}

func readRooms() []room.Room {
	rooms := readJsonFile[RoomJson](pathToRoomsJson)
	if rooms == nil {
		return nil
	}
	var roomEntities []room.Room

	for _, roomJson := range rooms {
		room := room.NewRoom(roomJson.Id, roomJson.Name, roomJson.Type, roomJson.UserIds)
		roomEntities = append(roomEntities, *room)
	}
	return roomEntities
}

func seedMessages(ctx context.Context, client *mongo.Client) {
	messageRepository := message.NewMessageRepository(client)
	messages := readMessages()
	fmt.Println("Deleting current messages...")
	messageRepository.DeleteAllMessages(ctx)
	fmt.Println("Adding new messages...")
	for _, message := range messages {
		messageRepository.AddMessage(ctx, &message)
	}
}

func readMessages() []message.Message {
	messages := readJsonFile[MessageJson](pathToMessagesJson)
	if messages == nil {
		return nil
	}
	var messageEntities []message.Message

	for _, messageJson := range messages {
		message := message.NewMessage(
			messageJson.Id,
			messageJson.Text,
			messageJson.CreatorID,
			messageJson.ReceiverID,
			messageJson.RoomID,
			messageJson.CreatedAt,
		)
		messageEntities = append(messageEntities, *message)
	}

	return messageEntities
}

func readJsonFile[T any](path string) []T {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer file.Close()

	items := []T{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&items)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return items
}
