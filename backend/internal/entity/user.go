package entity

import (
	"database/sql"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
)

type User struct {
	ID        int64          `json:"id" db:"id"`
	Login     string         `json:"login" db:"login"`
	Password  string         `json:"password" db:"password"`
	HotelID   sql.NullInt64  `json:"hotel_id" db:"hotel_id"`
	UserType  types.UserType `json:"user_type" db:"user_type"`
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" db:"updated_at"`
}
