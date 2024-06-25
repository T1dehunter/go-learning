package user

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	IsBanned bool
	RoomID   *int
}

func NewUser(id int, name string, email string, password string, isBanned bool, roomID *int) *User {
	return &User{Id: id, Name: name, Email: email, Password: password, IsBanned: isBanned, RoomID: roomID}
}
