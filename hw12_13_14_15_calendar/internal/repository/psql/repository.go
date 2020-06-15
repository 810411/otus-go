package psql

import (
	"context"
	"database/sql"
	"time"

	"github.com/jackc/pgconn"
	// Import Postgres sql driver
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/810411/otus-go/hw12_13_14_15_calendar/internal/repository"
)

type Repo struct {
	db *sql.DB
}

func New() *Repo {
	return new(Repo)
}

func (r *Repo) Connect(ctx context.Context, dsn string) (err error) {
	r.db, err = sql.Open("pgx", dsn)
	if err != nil {
		return
	}
	r.db.Stats()
	return r.db.PingContext(ctx)
}

func (r *Repo) Close() error {
	return r.db.Close()
}

func (r *Repo) Create(ctx context.Context, event repository.Event) (repository.Event, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return event, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO events (title, datetime, duration, description, owner_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		event.Title,
		event.Datetime.Format("2006-01-02 15:04:00 -0700"),
		event.Duration,
		event.Description,
		event.OwnerID,
	).Scan(&event.ID)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return event, repository.ErrTimeBusy
		}
		return event, err
	}

	if err = tx.Commit(); err != nil {
		return event, err
	}

	return event, nil
}

func (r *Repo) Update(ctx context.Context, event repository.Event) (repository.Event, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return event, err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	res, err := tx.ExecContext(
		ctx,
		`UPDATE events SET title = $1, datetime = $2, duration = $3, description = $4, updated_at = $5 WHERE id = $6`,
		event.Title,
		event.Datetime.Format("2006-01-02 15:04:00 -0700"),
		event.Duration,
		event.Description,
		"now()",
		event.ID,
	)
	if err != nil {
		pgErr, ok := err.(*pgconn.PgError)
		if ok && pgErr.Code == "23505" {
			return event, repository.ErrTimeBusy
		}
		return event, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return event, err
	}
	if ra == 0 {
		return event, repository.ErrNotFound
	}

	if err = tx.Commit(); err != nil {
		return event, err
	}

	return event, nil
}

func (r *Repo) Delete(ctx context.Context, id repository.EventID) (repository.EventID, error) {
	res, err := r.db.ExecContext(ctx, `DELETE FROM events WHERE id = $1`, id)
	if err != nil {
		return id, err
	}

	ra, err := res.RowsAffected()
	if err != nil {
		return id, err
	}
	if ra == 0 {
		return id, repository.ErrNotFound
	}

	return id, nil
}

func (r *Repo) listOf(ctx context.Context, from time.Time, p repository.Period) ([]repository.Event, error) {
	var events []repository.Event
	from, to := repository.GetTimeRange(from, p)

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, title, datetime, duration, description, owner_id FROM events WHERE datetime >= $1 AND datetime < $2`,
		from.Format("2006-01-02"),
		to.Format("2006-01-02"),
	)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var (
			event       repository.Event
			duration    sql.NullInt64
			description sql.NullString
		)

		if err := rows.Scan(
			&event.ID,
			&event.Title,
			&event.Datetime,
			&duration,
			&description,
			&event.OwnerID,
		); err != nil {
			return nil, err
		}

		if duration.Valid {
			event.Duration = time.Duration(duration.Int64)
		}
		if description.Valid {
			event.Description = description.String
		}

		events = append(events, event)
	}

	return events, rows.Err()
}

func (r *Repo) ListOfDay(ctx context.Context, from time.Time) ([]repository.Event, error) {
	return r.listOf(ctx, from, repository.Day)
}

func (r *Repo) ListOfWeek(ctx context.Context, from time.Time) ([]repository.Event, error) {
	return r.listOf(ctx, from, repository.Week)
}

func (r *Repo) ListOfMonth(ctx context.Context, from time.Time) ([]repository.Event, error) {
	return r.listOf(ctx, from, repository.Month)
}
