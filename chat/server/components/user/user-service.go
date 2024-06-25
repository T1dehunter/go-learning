package user

import "context"

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) FindUserById(ctx context.Context, id int) *User {
	return userService.userRepository.FindUserById(ctx, id)
}
