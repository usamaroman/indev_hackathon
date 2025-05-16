package response

import "github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"

type GetAllRooms struct {
	Rooms []entity.RoomType `json:"rooms"`
}

type RoomInfo struct {
	Room     *entity.Room `json:"room"`
	IsActive bool         `json:"is_active"`
}

type Token struct {
	Token string `json:"token"`
}