package response

import "github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"

type GetAllRooms struct {
	Rooms []entity.RoomType `json:"rooms"`
}

type RoomInfo struct {
	Room   *entity.Room `json:"room"`
	Status string       `json:"status"`
}

type Token struct {
	Token string `json:"token"`
}
