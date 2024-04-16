package user

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (userRepository *UserRepository) FindUserById(id int) *User {
	users := make(map[int]*User)

	user1RoomID := 1
	users[1] = NewUser(1, "Alex", "alex@mail.com", "Test1234", false, &user1RoomID)

	user2RoomID := 1
	users[2] = NewUser(2, "Bob", "bob@mail.com", "Test1234", false, &user2RoomID)

	users[3] = NewUser(3, "Zed", "zed@mail.com", "Test1234", false, nil)

	user, ok := users[id]
	if !ok {
		return nil
	}

	return user
}
