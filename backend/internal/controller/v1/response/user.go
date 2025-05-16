package response

import (
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
)

type User struct {
	ID        int64     `json:"id"`
	Login     string    `json:"login"`
	HotelID   int64     `json:"business_id"`
	UserType  string    `json:"user_type"`
	CreatedAt time.Time `json:"created_at"`
}

type GetAllUsers struct {
	Users []*User `json:"users"`
}

func BuildUser(user entity.User) *User {
	return &User{
		ID:        user.ID,
		Login:     user.Login,
		HotelID:   user.HotelID.Int64,
		UserType:  user.UserType.String(),
		CreatedAt: user.CreatedAt,
	}
}

func BuildGetAllUsers(users []entity.User) *GetAllUsers {
	resUsers := make([]*User, 0, len(users))

	for _, u := range users {
		resUsers = append(resUsers, BuildUser(u))
	}

	return &GetAllUsers{
		Users: resUsers,
	}
}