package user

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	client *mongo.Client
}

type UserSchema struct {
	Id       int
	Name     string
	Email    string
	IsBanned bool
	Password string
	RoomID   int
}

func NewUserRepository(client *mongo.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (userRepository *UserRepository) FindUserById(ctx context.Context, id int) *User {
	var result UserSchema
	filter := bson.D{{"id", id}}
	collection := userRepository.client.Database("chat").Collection("users")
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == nil {
		return mapSchemaToUser(&result)
	}
	return nil
}

func (userRepository *UserRepository) AddUser(ctx context.Context, user *User) {
	collection := userRepository.client.Database("chat").Collection("users")
	data := map[string]interface{}{
		"id":       user.Id,
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
		"isBanned": user.IsBanned,
		"roomID":   user.RoomID,
	}
	_, err := collection.InsertOne(ctx, data)
	if err != nil {
		panic(err)
	}
}

func (userRepository *UserRepository) DeleteAllUsers(ctx context.Context) {
	userRepository.client.Database("chat").Collection("users").Drop(ctx)
}
