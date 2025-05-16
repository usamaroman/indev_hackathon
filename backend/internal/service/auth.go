package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/repoerrors"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"

	"github.com/golang-jwt/jwt"
)

type AccessTokenClaims struct {
	jwt.StandardClaims
	UserID   int64  `json:"user_id"`
	HotelID  int64  `json:"hotel_id"`
	UserType string `json:"user_type"`
}

type RefreshTokenClaims struct {
	jwt.StandardClaims
	UserID int64 `json:"user_id"`
}

type AuthService struct {
	log *slog.Logger

	signKey  string
	tokenTTL time.Duration
	userRepo repo.User
}

func NewAuthService(log *slog.Logger, userRepo repo.User, signKey string, tokenTTL time.Duration) *AuthService {
	log = log.With(slog.String("component", "auth service"))

	return &AuthService{
		log:      log,
		signKey:  signKey,
		tokenTTL: tokenTTL,
		userRepo: userRepo,
	}
}

func (s *AuthService) GenerateTokens(ctx context.Context, userID int64) (string, string, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repoerrors.ErrNotFound) {
			return "", "", ErrUserNotFound
		}

		s.log.Error("failed to get user by id", slog.Int64("userID", userID), logger.Error(err))
		return "", "", err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &AccessTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID:   user.ID,
		UserType: user.UserType.String(),
		HotelID:  user.HotelID.Int64,
	})

	accessTokenString, err := accessToken.SignedString([]byte(s.signKey))
	if err != nil {
		s.log.Error("failed to sign token", logger.Error(err))

		return "", "", ErrSignToken
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &RefreshTokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(s.signKey))
	if err != nil {
		s.log.Error("failed to sign token", logger.Error(err))

		return "", "", ErrSignToken
	}

	return accessTokenString, refreshTokenString, nil
}

func (s *AuthService) ParseAccessToken(accessToken string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		s.log.Error("failed to parse token", logger.Error(err))
		return nil, ErrParseToken
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok {
		s.log.Error("failed to parse token", logger.Error(err))
		return nil, ErrParseToken
	}

	return claims, nil
}

func (s *AuthService) ParseRefreshToken(accessToken string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})
	if err != nil {
		s.log.Error("failed to parse token", logger.Error(err))
		return nil, ErrParseToken
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok {
		s.log.Error("failed to parse token", logger.Error(err))
		return nil, ErrParseToken
	}

	return claims, nil
}
