package app

import (
	"context"
	"time"

	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
	userID  int
}

type Logger interface {
	Debug(msg string, params ...any)
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, event storage.Event) error
	GetEvent(ctx context.Context, eventID string) (storage.Event, error)
	GetEventsListByDates(ctx context.Context, from *time.Time, to *time.Time) []storage.Event
	GetEventsForNotify(ctx context.Context, notifyDate string) []storage.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logger:  logger,
		storage: storage,
	}
}

func (a *App) CreateEvent(
	ctx context.Context,
	id, title string,
	startDate time.Time,
	endDate time.Time,
	notifyBefore time.Duration,
) error {
	return a.storage.CreateEvent(
		ctx,
		storage.Event{
			ID:           id,
			Title:        title,
			StartDate:    startDate,
			EndDate:      endDate,
			CreatorID:    a.userID,
			NotifyBefore: notifyBefore,
		},
	)
}

func (a *App) UpdateEvent(
	ctx context.Context,
	id, title string,
	startDate time.Time,
	endDate time.Time,
	notifyBefore time.Duration,
) error {
	return a.storage.UpdateEvent(
		ctx,
		storage.Event{
			ID:           id,
			Title:        title,
			StartDate:    startDate,
			EndDate:      endDate,
			CreatorID:    a.userID,
			NotifyBefore: notifyBefore,
		},
	)
}

func (a *App) GetEvent(ctx context.Context, id string) (storage.Event, error) {
	return a.storage.GetEvent(ctx, id)
}

func (a *App) GetEventsListByDates(ctx context.Context, from *time.Time, to *time.Time) []storage.Event {
	return a.storage.GetEventsListByDates(ctx, from, to)
}

func (a *App) GetEventsForNotify(ctx context.Context, notifyDate string) []storage.Event {
	return a.storage.GetEventsForNotify(ctx, notifyDate)
}
