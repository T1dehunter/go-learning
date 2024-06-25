package room

type RoomType string

const (
	Direct RoomType = "direct"
	Group  RoomType = "group"
)

type Room struct {
	Id      int
	Name    string
	Type    string
	UserIds []int
}

func NewRoom(id int, name string, roomType string, userIds []int) *Room {
	return &Room{Id: id, Name: name, Type: roomType, UserIds: userIds}
}

func (room *Room) JoinUser(userId int) {
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

func (room *Room) IsHasUser(userId int) bool {
	return contains(room.UserIds, userId)
}

func contains(arr []int, target int) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}
