package user

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
	IsBanned bool
}

func NewUser(id int, name string, email string, password string, isBanned bool) *User {
	return &User{Id: id, Name: name, Email: email, Password: password, IsBanned: isBanned}
}
