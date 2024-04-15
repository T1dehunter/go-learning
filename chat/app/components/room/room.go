package room

type Room struct {
	Id      int
	Name    string
	UserIds []int
}

func NewRoom(id int, name string, userIds []int) *Room {
	return &Room{Id: id, Name: name, UserIds: userIds}
}

func (room *Room) joinUser(userId int) {
	isContainsUser := contains(room.UserIds, userId)
	if isContainsUser {
		return
	}
	room.UserIds = append(room.UserIds, userId)
}

func (room *Room) leaveUser(userId int) {
	for userIdx, userID := range room.UserIds {
		if userID == userId {
			room.UserIds = append(room.UserIds[:userIdx], room.UserIds[userIdx+1:]...)
			return
		}
	}
}

func contains(arr []int, target int) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}
