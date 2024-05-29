package sqlstorage

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib" // postgres driver
	"github.com/jmoiron/sqlx"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

var ErrDBNotConnected = errors.New("not connected to database")

type SQLStorage struct {
	dsn string
	mu  sync.Mutex
	db  *sqlx.DB
}

type StorageEvent struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	Description  string        `db:"description"`
	StartDate    time.Time     `db:"start_dt"`
	EndDate      time.Time     `db:"end_dt"`
	CreatorID    int           `db:"creator_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

func New(dsn string) *SQLStorage {
	return &SQLStorage{
		dsn: dsn,
	}
}

func (s *SQLStorage) Connect(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.db != nil {
		// already connected
		return nil
	}

	db, err := sqlx.ConnectContext(ctx, "pgx", s.dsn)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *SQLStorage) Close(_ context.Context) error {
	err := s.db.Close()
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	s.db = nil
	return nil
}

func (s *SQLStorage) CreateEvent(ctx context.Context, event storage.Event) error {
	if s.db == nil {
		return ErrDBNotConnected
	}

	query := `INSERT INTO public.events (id, creator_id, title, description, start_dt, end_dt, notify_before)
			   VALUES (:id, :creator_id, :title, :description, :start_dt, :end_dt, :notify_before)`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":            event.ID,
		"creator_id":    event.CreatorID,
		"title":         event.Title,
		"description":   event.Description,
		"start_dt":      event.StartDate,
		"end_dt":        event.EndDate,
		"notify_before": event.NotifyBefore,
	})

	var e *pgconn.PgError
	if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
		return storage.ErrCreateEventIDExists
	}

	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) UpdateEvent(ctx context.Context, event storage.Event) error {
	if s.db == nil {
		return ErrDBNotConnected
	}

	query := `UPDATE public.events SET
			   creator_id = :creator_id, 
			   title = :title, 
			   description = :description, 
			   start_dt = :start_dt, 
			   end_dt = :end_dt, 
			   notify_before = :notify_before
			WHERE id = :id`

	_, err := s.db.NamedExecContext(ctx, query, map[string]interface{}{
		"id":            event.ID,
		"creator_id":    event.CreatorID,
		"title":         event.Title,
		"description":   event.Description,
		"start_dt":      event.StartDate,
		"end_dt":        event.EndDate,
		"notify_before": event.NotifyBefore,
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLStorage) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	if s.db == nil {
		return storage.Event{}, ErrDBNotConnected
	}

	query := "SELECT id, creator_id, title, description, start_dt, end_dt, notify_before FROM public.events WHERE id = :id"
	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"id": eventID,
	})
	if err != nil {
		return storage.Event{}, err
	}

	var event StorageEvent
	hasRows := rows.Next()
	if !hasRows {
		return storage.Event{}, storage.ErrReadEventNotExists
	}

	err = rows.StructScan(&event)
	if err != nil {
		return storage.Event{}, err
	}

	return storage.Event{
		ID:           event.ID,
		CreatorID:    event.CreatorID,
		Title:        event.Title,
		Description:  event.Description,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		NotifyBefore: event.NotifyBefore,
	}, nil
}

func (s *SQLStorage) GetEventsListByDates(ctx context.Context, from *time.Time, to *time.Time) []storage.Event {
	if s.db == nil {
		return []storage.Event{}
	}

	params := map[string]interface{}{}

	query := `SELECT id, creator_id, title, description, start_dt, end_dt, notify_before 
			  FROM public.events`

	if from != nil || to != nil {
		query += " WHERE 1=1"
		if from != nil {
			query += " AND start_dt >= :from"
			params["from"] = from
		}
		if to != nil {
			query += " AND end_dt <= :to"
			params["to"] = to
		}
	}

	rows, err := s.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return []storage.Event{}
	}

	events := []storage.Event{}
	for rows.Next() {
		var event StorageEvent
		err = rows.StructScan(&event)

		if err != nil {
			continue
		}
		events = append(events, storage.Event{
			ID:           event.ID,
			CreatorID:    event.CreatorID,
			Title:        event.Title,
			Description:  event.Description,
			StartDate:    event.StartDate,
			EndDate:      event.EndDate,
			NotifyBefore: event.NotifyBefore,
		})
	}
	return events
}

func (s *SQLStorage) GetEventsForNotify(ctx context.Context, notifyDate string) []storage.Event {
	if s.db == nil {
		return []storage.Event{}
	}

	query := `SELECT id, creator_id, title, description, start_dt, end_dt, notify_before 
			  FROM public.events
			  where cast((end_dt - cast(CONCAT(notify_before/1000000, ' milliseconds') as interval)) as date) = :notify_date`

	rows, err := s.db.NamedQueryContext(ctx, query, map[string]interface{}{
		"notify_date": notifyDate,
	})
	if err != nil {
		return []storage.Event{}
	}

	events := []storage.Event{}
	for rows.Next() {
		var event StorageEvent
		err = rows.StructScan(&event)

		if err != nil {
			continue
		}
		events = append(events, storage.Event{
			ID:           event.ID,
			CreatorID:    event.CreatorID,
			Title:        event.Title,
			Description:  event.Description,
			StartDate:    event.StartDate,
			EndDate:      event.EndDate,
			NotifyBefore: event.NotifyBefore,
		})
	}
	return events
}

func (s *SQLStorage) RemoveEvents(ctx context.Context) error {
	if s.db == nil {
		return ErrDBNotConnected
	}

	_, err := s.db.ExecContext(ctx, "DELETE FROM public.events")
	if err != nil {
		return err
	}

	return nil
}
