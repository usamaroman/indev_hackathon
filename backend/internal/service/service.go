package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo"
)

type AuthGenerateTokenInput struct {
	Login    string
	Password string
}

type Auth interface {
	GenerateTokens(ctx context.Context, userID int64) (string, string, error)
	ParseAccessToken(accessToken string) (*AccessTokenClaims, error)
	ParseRefreshToken(refreshToken string) (*RefreshTokenClaims, error)
}

type User interface {
	GetByID(ctx context.Context, id int64) (*entity.User, error)
	GetByLogin(ctx context.Context, login string) (*entity.User, error)
}

type CreateReservationInput struct {
	StartDate time.Time
	EndDate   time.Time
	RoomType  int64
	UserID    int64
}

type Hotel interface {
	GetAvailableRooms(ctx context.Context, start, end time.Time) ([]entity.RoomType, error)
	GetRoomByID(ctx context.Context, id string) (*entity.Room, error)
	CreateReservation(ctx context.Context, input *CreateReservationInput) error
	RoomHasReservations(ctx context.Context, id string) (bool, error)
	UpdateReservationStatus(ctx context.Context, id string, status types.ReservationType) error
	GetUserCurrentReservation(ctx context.Context, userID int64) (*entity.Reservation, error)
	RoomReservationStatus(ctx context.Context, id string) (string, error)
	GetReservationsByStatus(ctx context.Context, status string) ([]entity.ReservationInfo, error)
}

type Dependencies struct {
	Log   *slog.Logger
	Repos *repo.Repositories

	SignKey  string
	TokenTTL time.Duration
}

type Services struct {
	Auth  Auth
	User  User
	Hotel Hotel
}

func NewServices(deps *Dependencies) *Services {
	services := &Services{
		Auth:  NewAuthService(deps.Log, deps.Repos.User, deps.SignKey, deps.TokenTTL),
		User:  NewUserService(deps.Log, deps.Repos.User),
		Hotel: NewHotelService(deps.Log, deps.Repos.Hotel, deps.Repos.User),
	}

	return services
}
