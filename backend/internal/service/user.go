package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/repoerrors"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"
)

type UserService struct {
	log *slog.Logger

	userRepo repo.User
}

func NewUserService(log *slog.Logger, userRepo repo.User) *UserService {
	log = log.With(slog.String("component", "user service"))

	return &UserService{
		log:      log,
		userRepo: userRepo,
	}
}

func (s *UserService) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return &entity.User{}, ErrUserNotFound
		}
		s.log.Error("failed to get user by id in database", logger.Error(err))
		return &entity.User{}, err
	}
	return user, nil
}

func (s *UserService) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	user, err := s.userRepo.GetByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return &entity.User{}, ErrUserNotFound
		}
		s.log.Error("failed to get user by id in database", logger.Error(err))
		return &entity.User{}, err
	}
	return user, nil
}
