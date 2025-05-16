package request

type GetRooms struct {
	StartDate string `json:"start_date" validate:"required" example:"01.12.2025"`
	EndDate   string `json:"end_date" validate:"required" example:"29.05.2026"`
}

type ReserveRoom struct {
	RoomType  int64  `json:"room_type" validate:"required" example:"1"`
	StartDate string `json:"start_date" validate:"required" example:"01.12.2025"`
	EndDate   string `json:"end_date" validate:"required" example:"29.05.2026"`
}
