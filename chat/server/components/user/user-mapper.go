package user

func mapSchemaToUser(userSchema *UserSchema) *User {
	return NewUser(userSchema.Id, userSchema.Name, userSchema.Email, userSchema.Password, userSchema.IsBanned, &userSchema.RoomID)
}
