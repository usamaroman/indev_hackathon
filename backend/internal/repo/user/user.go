package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/codes"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/repoerrors"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/postgresql"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repo struct {
	log *slog.Logger
	*postgresql.Postgres
}

func NewRepo(log *slog.Logger, pg *postgresql.Postgres) *Repo {
	return &Repo{
		log:      log,
		Postgres: pg,
	}
}

func (r *Repo) GetByLogin(ctx context.Context, login string) (*entity.User, error) {
	q := "SELECT * FROM users WHERE login = $1"

	r.log.Debug("get user by login query", slog.String("query", q))

	var user entity.User
	err := r.Pool.QueryRow(ctx, q, login).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.UserType,
		&user.HotelID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		var pgErr *pgconn.PgError
		if ok := errors.As(err, &pgErr); ok {
			if pgErr.Code == codes.UniqueConstraintCode {
				return nil, repoerrors.ErrAlreadyExists
			}
		}

		return nil, err
	}

	return &user, nil
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	q := "SELECT * FROM users WHERE id = $1"

	r.log.Debug("get user by id query", slog.String("query", q))

	var user entity.User
	err := r.Pool.QueryRow(ctx, q, id).Scan(
		&user.ID,
		&user.Login,
		&user.Password,
		&user.UserType,
		&user.HotelID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}
