package entity

import (
	"time"
)

type Hotel struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type RoomType struct {
	ID       int64  `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Capacity int64  `json:"capacity" db:"capacity"`
}

type Room struct {
	Number     string `json:"room_number" db:"room_number"`
	HotelID    int64  `json:"hotel_id" db:"hotel_id"`
	RoomTypeID int64  `json:"room_type_id" db:"room_type_id"`
	Floor      int64  `json:"floor" db:"floor"`
}

type Reservation struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	CheckIn   time.Time `json:"check_in" db:"check_in"`
	CheckOut  time.Time `json:"check_out" db:"check_out"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ReservationInfo struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoomID    string    `json:"room_id" db:"room_id"`
	CheckIn   time.Time `json:"check_in" db:"check_in"`
	CheckOut  time.Time `json:"check_out" db:"check_out"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Login     string    `json:"login" db:"login"`
}
