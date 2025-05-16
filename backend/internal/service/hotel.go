package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo"
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

	user, err := s.userRepo.GetByID(ctx, input.UserID)
	if err != nil {
		return err
	}

	if err = s.hotelRepo.CreateReservation(ctx, &entity.Reservation{
		RoomID:    roomID,
		GuestName: user.Login,
		CheckIn:   input.StartDate,
		CheckOut:  input.EndDate,
	}); err != nil {
		return err
	}

	return nil
}
