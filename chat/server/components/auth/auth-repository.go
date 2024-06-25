package auth

import (
	"chat/server/components/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepository struct {
	userRepository *user.UserRepository
}

func NewAuthRepository(client *mongo.Client) *AuthRepository {
	return &AuthRepository{userRepository: user.NewUserRepository(client)}
}

func (authRepo *AuthRepository) FindUserByID(id int) *user.User {
	return nil
}
