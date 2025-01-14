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

func (userService *UserService) FindAllUsersByIds(ctx context.Context, ids []int) []*User {
	return userService.userRepository.FindAllUsersByIds(ctx, ids)
}

func (userService *UserService) FindUserByAccessToken(ctx context.Context, accessToken string) *User {
	return userService.userRepository.FindUserByAccessToken(ctx, accessToken)
}

func (userService *UserService) FindUserByNameAndPassword(ctx context.Context, name string, password string) *User {
	return userService.userRepository.FindUserByNameAndPassword(ctx, name, password)
}
