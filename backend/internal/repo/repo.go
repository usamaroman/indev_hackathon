package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/hotel"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/user"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/postgresql"
)

type User interface {
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
}

type Hotel interface {
	GetAvailableRooms(ctx context.Context, start, end time.Time) ([]entity.RoomType, error)
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error)
	CreateReservation(ctx context.Context, rsv *entity.Reservation) error
	GetAvailableRoomByType(ctx context.Context, roomType int64, start, end time.Time) (string, error)
	RoomHasReservations(ctx context.Context, id string) (bool, error)
	GetUserCurrentReservation(ctx context.Context, userID int64) (*entity.Reservation, error)
	UpdateReservationStatus(ctx context.Context, id string, status types.ReservationType) error
	RoomReservationStatus(ctx context.Context, id string) (string, error)
	GetReservationsByStatus(ctx context.Context, status string) ([]entity.ReservationInfo, error)
}

type Repositories struct {
	User
	Hotel
}

func NewRepositories(log *slog.Logger, pg *postgresql.Postgres) *Repositories {
	return &Repositories{
		User:  user.NewRepo(log, pg),
		Hotel: hotel.NewRepo(log, pg),
	}
}
