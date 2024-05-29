package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

var testDSN = "postgres://calendar:calendar@localhost:5432/calendar"

func TestStorageCreate(t *testing.T) {
	store := New(testDSN)

	newEvent := storage.Event{
		ID:           uuid.NewString(),
		CreatorID:    1,
		Title:        "Test",
		Description:  "Test description",
		StartDate:    time.Now().UTC(),
		EndDate:      time.Now().Add(time.Hour * 24 * 5).UTC(),
		NotifyBefore: time.Hour * 24 * 1,
	}
	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	err = store.CreateEvent(ctx, newEvent)
	require.Nil(t, err)

	savedEvent, err := store.GetEvent(ctx, newEvent.ID)
	require.Nil(t, err)
	require.Equal(t, newEvent.ID, savedEvent.ID)
	require.Equal(t, newEvent.CreatorID, savedEvent.CreatorID)
	require.Equal(t, newEvent.Title, savedEvent.Title)
	require.Equal(t, newEvent.Description, savedEvent.Description)
	require.Equal(t, newEvent.StartDate.Format("2006-01-02 15:04:05"), savedEvent.StartDate.Format("2006-01-02 15:04:05"))
	require.Equal(t, newEvent.EndDate.Format("2006-01-02 15:04:05"), savedEvent.EndDate.Format("2006-01-02 15:04:05"))
	require.Equal(t, newEvent.NotifyBefore, savedEvent.NotifyBefore)

	err = store.CreateEvent(ctx, newEvent)
	require.Equal(t, storage.ErrCreateEventIDExists, err)
}

func TestStorageUpdate(t *testing.T) {
	store := New(testDSN)

	startDate := time.Now()
	newEvent := storage.Event{
		ID:           uuid.NewString(),
		CreatorID:    1,
		Title:        "Test",
		Description:  "Test Description",
		StartDate:    startDate,
		EndDate:      startDate.Add(time.Duration(time.Duration.Hours(time.Hour * 24 * 6))),
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	}
	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	err = store.CreateEvent(ctx, newEvent)
	if err != nil {
		t.Fail()
	}
	newEvent.Title = "Test Test"
	newEvent.Description = "Test Description Test"
	err = store.UpdateEvent(ctx, newEvent)
	require.Nil(t, err)

	savedEvent, err := store.GetEvent(ctx, newEvent.ID)
	if err != nil {
		t.Fail()
	}

	require.Equal(t, "Test Test", savedEvent.Title)
	require.Equal(t, "Test Description Test", savedEvent.Description)
}

func TestStorageGet(t *testing.T) {
	store := New(testDSN)

	events := []storage.Event{
		{
			ID:    uuid.NewString(),
			Title: "Test",
		},
		{
			ID:    uuid.NewString(),
			Title: "Test 2",
		},
		{
			ID:    uuid.NewString(),
			Title: "Test 3",
		},
	}

	ctx := context.Background()
	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	for i := range events {
		err = store.CreateEvent(ctx, events[i])
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	for i := range events {
		event, err := store.GetEvent(ctx, events[i].ID)
		require.Nil(t, err)
		require.Equal(t, event, events[i])
	}

	_, err = store.GetEvent(ctx, uuid.NewString())
	require.Equal(t, storage.ErrReadEventNotExists, err)
}

func TestStorageGetEventsListByDates(t *testing.T) {
	store := New(testDSN)

	ctx := context.Background()
	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()
	uuid3 := uuid.NewString()

	startDate, _ := time.Parse("2006-01-02", "2024-05-26")
	endDate, _ := time.Parse("2006-01-02", "2024-06-26")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           uuid1,
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
		ID:           uuid2,
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
		ID:           uuid3,
		Title:        "Test 3",
		Description:  "Test Description 3",
		StartDate:    startDate,
		EndDate:      endDate,
		NotifyBefore: time.Duration(time.Duration.Hours(time.Hour * 24 * 1)),
	})
	if err != nil {
		t.Fatal(err)
	}

	fromDate, _ := time.Parse("2006-01-02", "2024-04-17")
	toDate, _ := time.Parse("2006-01-02", "2024-06-26")
	events := store.GetEventsListByDates(ctx, &fromDate, &toDate)

	require.Equal(t, 2, len(events))
	require.Equal(t, uuid1, events[0].ID)
	require.Equal(t, uuid2, events[1].ID)

	events = store.GetEventsListByDates(ctx, nil, nil)
	require.Equal(t, 3, len(events))
}

func TestStorageGetEventsForNotify(t *testing.T) {
	store := New(testDSN)

	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()
	uuid3 := uuid.NewString()

	startDate, _ := time.Parse("2006-01-02", "2024-05-26")
	endDate, _ := time.Parse("2006-01-02 15:04:05", "2024-06-26 23:59:59")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           uuid1,
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
		ID:           uuid2,
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
		ID:           uuid3,
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
	require.Equal(t, uuid1, events[0].ID)

	events = store.GetEventsForNotify(ctx, "2024-05-16")
	require.Equal(t, 1, len(events))
	require.Equal(t, uuid2, events[0].ID)

	events = store.GetEventsForNotify(ctx, "2024-07-19")
	require.Equal(t, 1, len(events))
	require.Equal(t, uuid3, events[0].ID)
}
