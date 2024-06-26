package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

type inMemoryEvent struct {
	ID           string
	Title        string
	Description  string
	StartDate    time.Time
	EndDate      time.Time
	CreatorID    int
	NotifyBefore time.Duration
}

type InMemoryStorage struct {
	mu   sync.RWMutex
	data map[string]inMemoryEvent
}

func New() *InMemoryStorage {
	return &InMemoryStorage{
		data: map[string]inMemoryEvent{},
	}
}

func (s *InMemoryStorage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[event.ID]
	if ok {
		return storage.ErrCreateEventIDExists
	}

	s.data[event.ID] = buildInMemoryEvent(event)

	return nil
}

func (s *InMemoryStorage) UpdateEvent(_ context.Context, eventID string, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	savedEvent, ok := s.data[eventID]
	if !ok {
		return storage.ErrUpdateEventIDNotExists
	}

	savedEvent = patchEventData(savedEvent, event)

	s.data[eventID] = savedEvent
	return nil
}

func (s *InMemoryStorage) DeleteEvent(_ context.Context, eventID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.data[eventID]
	if !ok {
		return nil
	}

	delete(s.data, eventID)

	return nil
}

func (s *InMemoryStorage) GetEvent(_ context.Context, eventID string) (storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	savedEvent, ok := s.data[eventID]
	if !ok {
		return storage.Event{}, storage.ErrReadEventNotExists
	}

	return buildStorageEvent(savedEvent), nil
}

func (s *InMemoryStorage) GetEventsListByDates(_ context.Context, from *time.Time, to *time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := []storage.Event{}

	for _, event := range s.data {
		if from != nil && event.EndDate.Before(*from) {
			continue
		}

		if to != nil && event.StartDate.After(*to) {
			continue
		}

		events = append(events, buildStorageEvent(event))
	}

	return events
}

func (s *InMemoryStorage) GetEventsForNotify(_ context.Context, notifyDate string) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := []storage.Event{}

	for _, event := range s.data {
		dateForNotify := event.EndDate.Add(-event.NotifyBefore)
		dateForNotifyStr := dateForNotify.Format("2006-01-02")
		if dateForNotifyStr != notifyDate {
			continue
		}

		events = append(events, buildStorageEvent(event))
	}

	return events
}

func (s *InMemoryStorage) GetEventsOnDate(_ context.Context, date time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := []storage.Event{}

	for _, event := range s.data {
		if !event.StartDate.Equal(date) {
			continue
		}

		events = append(events, buildStorageEvent(event))
	}

	return events
}

func (s *InMemoryStorage) GetEventsOnWeek(_ context.Context, weekStartDate time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	weekEndDate := weekStartDate.Add(time.Hour * 24 * 6)

	events := []storage.Event{}

	for _, event := range s.data {
		if (event.StartDate.Equal(weekStartDate) || event.StartDate.After(weekStartDate)) &&
			(event.EndDate.Equal(weekEndDate) || event.EndDate.Before(weekEndDate)) {
			events = append(events, buildStorageEvent(event))
		}
	}

	return events
}

func (s *InMemoryStorage) GetEventsOnMonth(_ context.Context, monthStartDate time.Time) []storage.Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	monthEndDate := monthStartDate.AddDate(0, 1, -1)

	events := []storage.Event{}

	for _, event := range s.data {
		if (event.StartDate.Equal(monthStartDate) || event.StartDate.After(monthStartDate)) &&
			(event.EndDate.Equal(monthEndDate) || event.EndDate.Before(monthEndDate)) {
			events = append(events, buildStorageEvent(event))
		}
	}

	return events
}

func buildStorageEvent(event inMemoryEvent) storage.Event {
	return storage.Event{
		ID:           event.ID,
		Title:        event.Title,
		Description:  event.Description,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		CreatorID:    event.CreatorID,
		NotifyBefore: event.NotifyBefore,
	}
}

func buildInMemoryEvent(event storage.Event) inMemoryEvent {
	return inMemoryEvent{
		ID:           event.ID,
		Title:        event.Title,
		Description:  event.Description,
		StartDate:    event.StartDate,
		EndDate:      event.EndDate,
		CreatorID:    event.CreatorID,
		NotifyBefore: event.NotifyBefore,
	}
}

func patchEventData(savedEvent inMemoryEvent, event storage.Event) inMemoryEvent {
	savedEvent.Title = event.Title
	savedEvent.Description = event.Description
	savedEvent.StartDate = event.StartDate
	savedEvent.EndDate = event.EndDate
	savedEvent.CreatorID = event.CreatorID
	savedEvent.NotifyBefore = event.NotifyBefore

	return savedEvent
}
