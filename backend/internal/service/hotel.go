package service

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/repoerrors"
)

type HotelService struct {
	log *slog.Logger

	hotelRepo repo.Hotel
	userRepo  repo.User
}

func NewHotelService(log *slog.Logger, hotelRepo repo.Hotel, userRepo repo.User) *HotelService {
	log = log.With(slog.String("component", "business service"))

	return &HotelService{
		log:       log,
		hotelRepo: hotelRepo,
		userRepo:  userRepo,
	}
}

func (s *HotelService) GetAvailableRooms(ctx context.Context, start, end time.Time) ([]entity.RoomType, error) {
	return s.hotelRepo.GetAvailableRooms(ctx, start, end)
}

func (s *HotelService) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	return s.hotelRepo.GetRoomByID(ctx, id)
}

func (s *HotelService) RoomHasReservations(ctx context.Context, id string) (bool, error) {
	return s.hotelRepo.RoomHasReservations(ctx, id)
}

func (s *HotelService) CreateReservation(ctx context.Context, input *CreateReservationInput) error {
	roomID, err := s.hotelRepo.GetAvailableRoomByType(ctx, input.RoomType, input.StartDate, input.EndDate)
	if err != nil {
		return err
	}

	if err = s.hotelRepo.CreateReservation(ctx, &entity.Reservation{
		RoomID:   roomID,
		UserID:   input.UserID,
		CheckIn:  input.StartDate,
		CheckOut: input.EndDate,
	}); err != nil {
		return err
	}

	return nil
}

func (s *HotelService) GetUserCurrentReservation(ctx context.Context, userID int64) (*entity.Reservation, error) {
	rsv, err := s.hotelRepo.GetUserCurrentReservation(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return nil, ErrReservationNotFound
		}

		return nil, err
	}

	return rsv, nil
}

func (s *HotelService) UpdateReservationStatus(ctx context.Context, id string, status types.ReservationType) error {
	return s.hotelRepo.UpdateReservationStatus(ctx, id, status)
}
