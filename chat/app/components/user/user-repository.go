package user

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (userRepository *UserRepository) FindUserById(id int) *User {
	users := make(map[int]*User)

	users[1] = NewUser(1, "Alex", "alex@mail.com", "Test1234", false)
	users[2] = NewUser(2, "Bob", "bob@mail.com", "Test1234", false)
	users[3] = NewUser(3, "Zed", "zed@mail.com", "Test1234", false)

	user, ok := users[id]
	if !ok {
		return nil
	}

	return user
}
