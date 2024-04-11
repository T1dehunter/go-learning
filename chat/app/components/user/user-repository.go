package user

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (userRepository *UserRepository) FindUserById(id int) *User {
	users := make(map[int]*User)

	users[1] = NewUser(1, "Alex", "alex@mail.com", "Test1234", false)

	user, ok := users[id]
	if !ok {
		return nil
	}

	return user
}
