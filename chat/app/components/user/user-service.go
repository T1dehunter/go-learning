package user

type UserService struct {
	userRepository *UserRepository
}

func NewUserService() *UserService {
	return &UserService{}
}

func (userService *UserService) FindUserById(id int) *User {
	return userService.userRepository.FindUserById(id)
}
