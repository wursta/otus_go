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

	defer store.RemoveEvents(ctx)

	err = store.CreateEvent(ctx, newEvent)
	require.Nil(t, err)

	savedEvent, err := store.GetEvent(ctx, newEvent.ID)
	require.Nil(t, err)
	require.Equal(t, newEvent.ID, savedEvent.ID)
	require.Equal(t, newEvent.CreatorID, savedEvent.CreatorID)
	require.Equal(t, newEvent.Title, savedEvent.Title)
	require.Equal(t, newEvent.Description, savedEvent.Description)
	require.Equal(t, newEvent.StartDate.Format(time.DateOnly), savedEvent.StartDate.Format(time.DateOnly))
	require.Equal(t, newEvent.EndDate.Format(time.DateOnly), savedEvent.EndDate.Format(time.DateOnly))
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

	defer store.RemoveEvents(ctx)

	err = store.CreateEvent(ctx, newEvent)
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

	defer store.RemoveEvents(ctx)

	for i := range events {
		err = store.CreateEvent(ctx, events[i])
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = store.GetEvent(ctx, events[1].ID)
	require.Nil(t, err)

	err = store.DeleteEvent(ctx, events[1].ID)
	require.Nil(t, err)

	_, err = store.GetEvent(ctx, events[1].ID)
	require.NotNil(t, err)

	err = store.DeleteEvent(ctx, events[1].ID)
	require.Nil(t, err)
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

	defer store.RemoveEvents(ctx)

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

	defer store.RemoveEvents(ctx)

	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()
	uuid3 := uuid.NewString()

	startDate, _ := time.Parse(time.DateOnly, "2024-05-26")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-26")
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

	startDate, _ = time.Parse(time.DateOnly, "2024-04-18")
	endDate, _ = time.Parse(time.DateOnly, "2024-05-18")
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

	startDate, _ = time.Parse(time.DateOnly, "2024-04-20")
	endDate, _ = time.Parse(time.DateOnly, "2024-07-22")
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

	fromDate, _ := time.Parse(time.DateOnly, "2024-04-17")
	toDate, _ := time.Parse(time.DateOnly, "2024-06-26")
	events := store.GetEventsListByDates(ctx, &fromDate, &toDate)

	require.Equal(t, 2, len(events))
	require.Equal(t, uuid1, events[0].ID)
	require.Equal(t, uuid2, events[1].ID)

	events = store.GetEventsListByDates(ctx, nil, nil)
	require.Equal(t, 3, len(events))
}

func TestGetEventsOnDate(t *testing.T) {
	store := New(testDSN)

	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	defer store.RemoveEvents(ctx)

	uuid1 := uuid.NewString()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
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

	events := store.GetEventsOnWeek(ctx, startDate)
	require.Equal(t, 1, len(events))
	require.Equal(t, uuid1, events[0].ID)
}

func TestGetEventsOnWeek(t *testing.T) {
	store := New(testDSN)

	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	defer store.RemoveEvents(ctx)

	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
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

	startDate, _ = time.Parse(time.DateOnly, "2024-06-05")
	endDate, _ = time.Parse(time.DateOnly, "2024-06-12")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           uuid2,
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
	require.Equal(t, uuid1, events[0].ID)
}

func TestGetEventsOnMonth(t *testing.T) {
	store := New(testDSN)

	ctx := context.Background()

	err := store.Connect(ctx)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer store.Close(ctx)

	defer store.RemoveEvents(ctx)

	uuid1 := uuid.NewString()
	uuid2 := uuid.NewString()
	uuid3 := uuid.NewString()

	startDate, _ := time.Parse(time.DateOnly, "2024-06-03")
	endDate, _ := time.Parse(time.DateOnly, "2024-06-05")
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

	startDate, _ = time.Parse(time.DateOnly, "2024-06-05")
	endDate, _ = time.Parse(time.DateOnly, "2024-06-12")
	err = store.CreateEvent(ctx, storage.Event{
		ID:           uuid2,
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
		ID:           uuid3,
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
	require.Equal(t, uuid1, events[0].ID)
	require.Equal(t, uuid2, events[1].ID)
}
