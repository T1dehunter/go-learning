package user

type UserService struct {
	userRepository *UserRepository
}

func NewUserService(userRepository *UserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) FindUserById(id int) *User {
	return userService.userRepository.FindUserById(id)
}
