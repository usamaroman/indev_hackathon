package hotel

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/entity/types"
	"github.com/usamaroman/demo_indev_hackathon/backend/internal/repo/repoerrors"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/logger"
	"github.com/usamaroman/demo_indev_hackathon/backend/pkg/postgresql"

	"github.com/jackc/pgx/v5"
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

func (r *Repo) GetAvailableRooms(ctx context.Context, start, end time.Time) ([]entity.RoomType, error) {
	q := `
SELECT rt.* 
FROM room_types rt 
WHERE rt.id IN (
    SELECT DISTINCT r.room_type_id
    FROM rooms r
    WHERE NOT EXISTS (
        SELECT 1
        FROM reservations res
        WHERE res.room_id = r.room_number
        AND res.status IN ('confirmed', 'checked_in')
        AND (res.check_in < $2 AND res.check_out > $1)
	)
)
`

	r.log.Debug("get available rooms query", slog.String("query", q))

	rows, err := r.Pool.Query(ctx, q, start, end)
	if err != nil {
		r.log.Error("failed to get rooms from database", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	rooms, err := pgx.CollectRows(rows, pgx.RowToStructByName[entity.RoomType])
	if err != nil {
		r.log.Error("failed to collect rows", logger.Error(err))
		return nil, err
	}

	return rooms, nil
}

func (r *Repo) GetRoomInfoByID(ctx context.Context, id string) (*entity.Room, error) {
	q := "SELECT * FROM rooms WHERE room_number = $1"

	r.log.Debug("get room query", slog.String("query", q), slog.String("id", id))

	rows, err := r.Pool.Query(ctx, q, id)
	if err != nil {
		r.log.Error("failed to get room from database", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	room, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Room])
	if err != nil {
		r.log.Error("failed to collect rows", logger.Error(err))
		return nil, err
	}

	return &room, nil
}

func (r *Repo) RoomHasReservations(ctx context.Context, id string) (bool, error) {
	q := "select count(*) from rooms r join reservations res on r.room_number = res.room_id where room_number like $1 and current_date between res.check_in and res.check_out"

	r.log.Debug("get if room has reservations query", slog.String("query", q), slog.String("id", id))

	var cnt int64

	if err := r.Pool.QueryRow(ctx, q, id).Scan(&cnt); err != nil {
		r.log.Error("failed to get room from database", logger.Error(err))
		return false, err
	}

	if cnt == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

func (r *Repo) GetRoomByID(ctx context.Context, id string) (*entity.Room, error) {
	q := "SELECT * FROM rooms WHERE room_number = $1"

	r.log.Debug("get room query", slog.String("query", q), slog.String("id", id))

	rows, err := r.Pool.Query(ctx, q, id)
	if err != nil {
		r.log.Error("failed to get room from database", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	room, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Room])
	if err != nil {
		r.log.Error("failed to collect rows", logger.Error(err))
		return nil, err
	}

	return &room, nil
}

func (r *Repo) GetAvailableRoomByType(ctx context.Context, roomType int64, start, end time.Time) (string, error) {
	q := `
SELECT r.room_number
FROM rooms r 
WHERE r.room_type_id = $1 
AND NOT EXISTS (
    SELECT 1 
    FROM reservations res 
    WHERE res.room_id = r.room_number 
    AND res.status IN ('confirmed', 'checked_in') 
    AND (res.check_in < $3 AND res.check_out > $2)
)
LIMIT 1;
`
	r.log.Debug("get available room by type query", slog.String("query", q), slog.Int64("room type", roomType))

	var id string
	if err := r.Pool.QueryRow(ctx, q, roomType, start, end).Scan(&id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", errors.New("no available rooms")
		}

		r.log.Error("failed to get available room by type from database", logger.Error(err))
		return "", err
	}

	return id, nil
}

func (r *Repo) CreateReservation(ctx context.Context, rsv *entity.Reservation) error {
	q := "INSERT INTO reservations (room_id, user_id, check_in, check_out) VALUES ($1, $2, $3, $4)"

	if _, err := r.Pool.Exec(ctx, q, rsv.RoomID, rsv.UserID, rsv.CheckIn, rsv.CheckOut); err != nil {
		r.log.Error("failed to create reservation", logger.Error(err))
		return err
	}

	return nil
}

func (r *Repo) GetUserCurrentReservation(ctx context.Context, userID int64) (*entity.Reservation, error) {
	q := "SELECT * FROM reservations WHERE user_id = $1 AND status = 'checked_in'"

	r.log.Debug("get user current reservation query", slog.String("query", q), slog.Int64("user id", userID))

	rows, err := r.Pool.Query(ctx, q, userID)
	if err != nil {
		r.log.Error("failed to get user current reservation from database", logger.Error(err))
		return nil, err
	}
	defer rows.Close()

	res, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[entity.Reservation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repoerrors.ErrNotFound
		}

		r.log.Error("failed to collect rows", logger.Error(err))
		return nil, err
	}

	return &res, nil
}

func (r *Repo) UpdateReservationStatus(ctx context.Context, id string, status types.ReservationType) error {
	q := "UPDATE reservations SET status = $1 WHERE id = $2"

	r.log.Debug("update reservation status query", slog.String("query", q), slog.String("id", id), slog.String("status", status.String()))

	if _, err := r.Pool.Exec(ctx, q, status, id); err != nil {
		r.log.Error("failed to update reservation status", logger.Error(err))
		return err
	}

	return nil
}
