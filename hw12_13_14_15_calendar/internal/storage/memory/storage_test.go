package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

func TestStorageCreate(t *testing.T) {
	store := New()

	newEvent := storage.Event{
		ID:    "1",
		Title: "Test",
	}
	ctx := context.Background()

	err := store.CreateEvent(ctx, newEvent)
	require.Nil(t, err)

	savedEvent, err := store.GetEvent(ctx, newEvent.ID)
	require.Nil(t, err)
	require.Equal(t, newEvent, savedEvent)

	err = store.CreateEvent(ctx, newEvent)
	require.Equal(t, storage.ErrCreateEventIDExists, err)
}

func TestStorageUpdate(t *testing.T) {
	store := New()

	startDate := time.Now()
	newEvent := storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      startDate.Add(time.Duration(time.Duration.Hours(time.Hour * 24 * 6))),
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	}
	ctx := context.Background()

	err := store.CreateEvent(ctx, newEvent)
	if err != nil {
		t.Fail()
	}
	newEvent.Title = "Test Test"
	newEvent.Description = "Test Description Test"
	err = store.UpdateEvent(ctx, newEvent.ID, newEvent)
	require.Nil(t, err)

	savedEvent, err := store.GetEvent(ctx, newEvent.ID)
	if err != nil {
		t.Fail()
	}

	require.Equal(t, "Test Test", savedEvent.Title)
	require.Equal(t, "Test Description Test", savedEvent.Description)
}

func TestStorageDelete(t *testing.T) {
	store := New()

	events := []storage.Event{
		{
			ID:    "1",
			Title: "Test",
		},
		{
			ID:    "2",
			Title: "Test 2",
		},
		{
			ID:    "3",
			Title: "Test 3",
		},
	}

	ctx := context.Background()

	for i := range events {
		err := store.CreateEvent(ctx, events[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err := store.GetEvent(ctx, events[1].ID)
	require.Nil(t, err)

	err = store.DeleteEvent(ctx, events[1].ID)
	require.Nil(t, err)

	_, err = store.GetEvent(ctx, events[1].ID)
	require.NotNil(t, err)

	err = store.DeleteEvent(ctx, events[1].ID)
	require.Nil(t, err)
}

func TestStorageGet(t *testing.T) {
	store := New()

	events := []storage.Event{
		{
			ID:    "1",
			Title: "Test",
		},
		{
			ID:    "2",
			Title: "Test 2",
		},
		{
			ID:    "3",
			Title: "Test 3",
		},
	}

	ctx := context.Background()
	for i := range events {
		err := store.CreateEvent(ctx, events[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	for i := range events {
		event, err := store.GetEvent(ctx, events[i].ID)
		require.Nil(t, err)
		require.Equal(t, event, events[i])
	}

	_, err := store.GetEvent(ctx, "4")
	require.Equal(t, storage.ErrReadEventNotExists, err)
}

func TestStorageGetEventsListByDates(t *testing.T) {
	store := New()

	ctx := context.Background()

	startDate, _ := time.Parse("2006-01-02", "2024-05-26")
	endDate, _ := time.Parse("2006-01-02", "2024-06-26")
	err := store.CreateEvent(ctx, storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse("2006-01-02", "2024-04-18")
	endDate, _ = time.Parse("2006-01-02", "2024-05-18")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "2",
		Title:        "Test 2",
		Description:  "Test Description 2",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse("2006-01-02", "2024-04-20")
	endDate, _ = time.Parse("2006-01-02", "2024-07-22")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "3",
		Title:        "Test 3",
		Description:  "Test Description 3",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	fromDate, _ := time.Parse("2006-01-02", "2024-05-20")
	toDate, _ := time.Parse("2006-01-02", "2024-05-28")
	events := store.GetEventsListByDates(ctx, &fromDate, &toDate)

	require.Equal(t, 2, len(events))
	require.Equal(t, "1", events[0].ID)
	require.Equal(t, "3", events[1].ID)

	events = store.GetEventsListByDates(ctx, nil, nil)
	require.Equal(t, 3, len(events))
}

func TestStorageGetEventsForNotify(t *testing.T) {
	store := New()

	ctx := context.Background()

	startDate, _ := time.Parse("2006-01-02", "2024-05-26")
	endDate, _ := time.Parse("2006-01-02 15:04:05", "2024-06-26 23:59:59")
	err := store.CreateEvent(ctx, storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse("2006-01-02", "2024-04-18")
	endDate, _ = time.Parse("2006-01-02 15:04:05", "2024-05-18 23:59:59")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "2",
		Title:        "Test 2",
		Description:  "Test Description 2",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 2,
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse("2006-01-02", "2024-04-20")
	endDate, _ = time.Parse("2006-01-02 15:04:05", "2024-07-22 23:59:59")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "3",
		Title:        "Test 3",
		Description:  "Test Description 3",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 3,
	})
	if err != nil {
		t.Fatal(err)
	}

	events := store.GetEventsForNotify(ctx, "2024-06-25")
	require.Equal(t, 1, len(events))
	require.Equal(t, "1", events[0].ID)

	events = store.GetEventsForNotify(ctx, "2024-05-16")
	require.Equal(t, 1, len(events))
	require.Equal(t, "2", events[0].ID)

	events = store.GetEventsForNotify(ctx, "2024-07-19")
	require.Equal(t, 1, len(events))
	require.Equal(t, "3", events[0].ID)
}

func TestGetEventsOnDate(t *testing.T) {
	store := New()

	ctx := context.Background()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
	err := store.CreateEvent(ctx, storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	events := store.GetEventsOnWeek(ctx, startDate)
	require.Equal(t, 1, len(events))
	require.Equal(t, "1", events[0].ID)
}

func TestGetEventsOnWeek(t *testing.T) {
	store := New()

	ctx := context.Background()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
	err := store.CreateEvent(ctx, storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse(time.DateOnly, "2024-06-05")
	endDate, _ = time.Parse(time.DateOnly, "2024-06-12")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "2",
		Title:        "Test 2",
		Description:  "Test Description 2",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	weekStartDate, _ := time.Parse(time.DateOnly, "2024-06-03")

	events := store.GetEventsOnWeek(ctx, weekStartDate)
	require.Equal(t, 1, len(events))
	require.Equal(t, "1", events[0].ID)
}

func TestGetEventsOnMonth(t *testing.T) {
	store := New()

	ctx := context.Background()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
	err := store.CreateEvent(ctx, storage.Event{
		ID:           "1",
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse(time.DateOnly, "2024-06-05")
	endDate, _ = time.Parse(time.DateOnly, "2024-06-12")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "2",
		Title:        "Test 2",
		Description:  "Test Description 2",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	startDate, _ = time.Parse(time.DateOnly, "2024-06-20")
	endDate, _ = time.Parse(time.DateOnly, "2024-07-05")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           "3",
		Title:        "Test 3",
		Description:  "Test Description 3",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Hour * 24 * 1,
	})
	if err != nil {
		t.Fatal(err)
	}

	weekStartDate, _ := time.Parse(time.DateOnly, "2024-06-01")

	events := store.GetEventsOnMonth(ctx, weekStartDate)
	require.Equal(t, 2, len(events))
}
