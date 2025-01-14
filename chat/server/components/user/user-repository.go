package user

import (
	"context"
	"fmt"
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

func (userRepository *UserRepository) FindAllUsersByIds(ctx context.Context, ids []int) []*User {
	var result UserSchema
	filter := bson.D{{"id", bson.D{{"$in", ids}}}}
	collection := userRepository.client.Database("chat").Collection("users")
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		panic(err)
	}
	defer cursor.Close(ctx)

	users := make([]*User, 0)
	for cursor.Next(ctx) {
		err := cursor.Decode(&result)
		if err != nil {
			panic(err)
		}
		users = append(users, mapSchemaToUser(&result))
	}

	fmt.Println("users", users)

	return users
}

func (userRepository *UserRepository) FindUserByAccessToken(ctx context.Context, accessToken string) *User {
	var result UserSchema
	// TEMP: for testing purposes, we are using the password field as the access token
	filter := bson.D{{"password", accessToken}}
	collection := userRepository.client.Database("chat").Collection("users")
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err == nil {
		return mapSchemaToUser(&result)
	}
	return nil
}

// TEMP: for testing purposes
func (userRepository *UserRepository) FindUserByNameAndPassword(ctx context.Context, name string, password string) *User {
	var result UserSchema
	filter := bson.D{{"name", name}, {"password", password}}

	fmt.Println("filter", filter)

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
