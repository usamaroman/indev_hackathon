package repo

import (
	"context"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/hotel"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/user"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/postgresql"
)

type User interface {
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	GetAllByBusinessID(ctx context.Context, businessID string) ([]entity.User, error)
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	DeleteByID(ctx context.Context, id string) error
}

type Hotel interface {
	GetAvailableRooms(ctx context.Context, start, end time.Time) ([]entity.RoomType, error)
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error)
	CreateReservation(ctx context.Context, rsv *entity.Reservation) error
	GetAvailableRoomByType(ctx context.Context, roomType int64, start, end time.Time) (string, error)
	RoomHasReservations(ctx context.Context, id string) (bool, error)
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
