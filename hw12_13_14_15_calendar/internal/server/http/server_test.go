package internalhttp

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/app"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/logger"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
)

func TestCreateHandler(t *testing.T) {
	var output bytes.Buffer
	logger, err := logger.New("DEBUG", &output)
	if err != nil {
		t.Fatal(err)
	}

	storage := memorystorage.New()
	app := app.New(logger, storage)
	timeout, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	server := NewServer(logger, app, "localhost", "8080", timeout)

	r := httptest.NewRequest("POST", "http://localhost:8080/event/create", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	server.CreateEventHandler(w, r)

	resp := w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	data := url.Values{}
	data.Set("id", "111")
	data.Set("title", "Test")
	data.Set("start_dt", "2025-06-01")
	data.Set("end_dt", "2025-07-01")
	data.Set("notify_before", "48h")

	r = httptest.NewRequest(
		"POST",
		"http://localhost:8080/event/create",
		strings.NewReader(data.Encode()),
	)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	server.CreateEventHandler(w, r)

	resp = w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	_, err = storage.GetEvent(context.Background(), "111")
	require.Nil(t, err)
}

func TestUpdateHandler(t *testing.T) {
	var output bytes.Buffer

	logger, err := logger.New("DEBUG", &output)
	if err != nil {
		t.Fatal(err)
	}

	store := memorystorage.New()

	app := app.New(logger, store)

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	server := NewServer(logger, app, "localhost", "8080", timeout)

	r := httptest.NewRequest("POST", "http://localhost:8080/event/update", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	server.UpdateEventHandler(w, r)

	resp := w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	startDt, err := time.Parse(time.DateOnly, "2025-06-01")
	if err != nil {
		t.Fatal(err)
	}

	endDt, err := time.Parse(time.DateOnly, "2025-07-01")
	if err != nil {
		t.Fatal(err)
	}

	store.CreateEvent(context.Background(), storage.Event{
		ID:           "111",
		Title:        "Test",
		StartDate:    startDt,
		EndDate:      endDt,
		NotifyBefore: time.Hour * 48,
	})

	data := url.Values{}
	data.Set("id", "111")
	data.Set("title", "Test Test")
	data.Set("start_dt", "2025-06-01")
	data.Set("end_dt", "2025-07-01")
	data.Set("notify_before", "48h")

	r = httptest.NewRequest(
		"POST",
		"http://localhost:8080/event/update",
		strings.NewReader(data.Encode()),
	)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	w = httptest.NewRecorder()
	server.UpdateEventHandler(w, r)

	resp = w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	event, err := store.GetEvent(context.Background(), "111")
	require.Nil(t, err)
	require.Equal(t, "Test Test", event.Title)
}

func TestGetHandler(t *testing.T) {
	var output bytes.Buffer

	logger, err := logger.New("DEBUG", &output)
	if err != nil {
		t.Fatal(err)
	}

	memStorage := memorystorage.New()

	app := app.New(logger, memStorage)

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	server := NewServer(logger, app, "localhost", "8080", timeout)
	r := httptest.NewRequest("POST", "http://localhost:8080/event/get", nil)

	w := httptest.NewRecorder()
	server.GetEventHandler(w, r)

	resp := w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusBadRequest, resp.StatusCode)

	r = httptest.NewRequest("POST", "http://localhost:8080/event/get?id=unknown", nil)

	w = httptest.NewRecorder()
	server.GetEventHandler(w, r)

	resp = w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	startDt, err := time.Parse(time.DateOnly, "2025-06-01")
	if err != nil {
		t.Fatal(err)
	}

	endDt, err := time.Parse(time.DateOnly, "2025-07-01")
	if err != nil {
		t.Fatal(err)
	}

	memStorage.CreateEvent(context.Background(), storage.Event{
		ID:           "111",
		Title:        "Test",
		StartDate:    startDt,
		EndDate:      endDt,
		CreatorID:    1,
		NotifyBefore: time.Hour * 48,
	})

	r = httptest.NewRequest(
		"POST",
		"http://localhost:8080/event/get?id=111",
		nil,
	)

	w = httptest.NewRecorder()
	server.GetEventHandler(w, r)

	resp = w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)
	buf := new(strings.Builder)

	_, err = io.Copy(buf, resp.Body)
	require.Nil(t, err)

	require.Equal(
		t,
		//nolint: all
		"{\"id\":\"111\",\"title\":\"Test\",\"description\":\"\",\"start_dt\":\"2025-06-01T00:00:00Z\",\"end_dt\":\"2025-07-01T00:00:00Z\",\"creator_id\":1,\"notify_before\":172800000000000}",
		buf.String(),
	)
}

func TestGetListByDatesHandler(t *testing.T) {
	var output bytes.Buffer

	logger, err := logger.New("DEBUG", &output)
	if err != nil {
		t.Fatal(err)
	}

	memStorage := memorystorage.New()

	app := app.New(logger, memStorage)

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	server := NewServer(logger, app, "localhost", "8080", timeout)

	startDt1, err := time.Parse(time.DateOnly, "2025-06-01")
	if err != nil {
		t.Fatal(err)
	}

	endDt1, err := time.Parse(time.DateOnly, "2025-06-10")
	if err != nil {
		t.Fatal(err)
	}

	startDt2, err := time.Parse(time.DateOnly, "2025-06-12")
	if err != nil {
		t.Fatal(err)
	}

	endDt2, err := time.Parse(time.DateOnly, "2025-06-15")
	if err != nil {
		t.Fatal(err)
	}

	memStorage.CreateEvent(context.Background(), storage.Event{
		ID:           "1",
		Title:        "Test",
		StartDate:    startDt1,
		EndDate:      endDt1,
		CreatorID:    1,
		NotifyBefore: time.Hour * 48,
	})

	memStorage.CreateEvent(context.Background(), storage.Event{
		ID:           "2",
		Title:        "Test 2",
		StartDate:    startDt2,
		EndDate:      endDt2,
		CreatorID:    1,
		NotifyBefore: time.Hour * 48,
	})

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/event/listByDates?from=2025-06-01&to=2025-06-10",
		nil,
	)

	w := httptest.NewRecorder()
	server.GetListByDatesHandler(w, r)

	resp := w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)

	require.Nil(t, err)
	require.Equal(
		t,
		//nolint: all
		"[{\"id\":\"1\",\"title\":\"Test\",\"description\":\"\",\"start_dt\":\"2025-06-01T00:00:00Z\",\"end_dt\":\"2025-06-10T00:00:00Z\",\"creator_id\":1,\"notify_before\":172800000000000}]",
		buf.String(),
	)
}

func TestGetListByNotifyDateHandler(t *testing.T) {
	var output bytes.Buffer

	logger, err := logger.New("DEBUG", &output)
	if err != nil {
		t.Fatal(err)
	}

	memStorage := memorystorage.New()

	app := app.New(logger, memStorage)

	timeout, err := time.ParseDuration("30s")
	if err != nil {
		t.Fatal(err)
	}

	server := NewServer(logger, app, "localhost", "8080", timeout)

	startDt1, err := time.Parse(time.DateOnly, "2025-06-01")
	if err != nil {
		t.Fatal(err)
	}

	endDt1, err := time.Parse(time.DateOnly, "2025-06-10")
	if err != nil {
		t.Fatal(err)
	}

	startDt2, err := time.Parse(time.DateOnly, "2025-06-12")
	if err != nil {
		t.Fatal(err)
	}

	endDt2, err := time.Parse(time.DateOnly, "2025-06-15")
	if err != nil {
		t.Fatal(err)
	}

	memStorage.CreateEvent(context.Background(), storage.Event{
		ID:           "1",
		Title:        "Test",
		StartDate:    startDt1,
		EndDate:      endDt1,
		CreatorID:    1,
		NotifyBefore: time.Hour * 48,
	})

	memStorage.CreateEvent(context.Background(), storage.Event{
		ID:           "2",
		Title:        "Test 2",
		StartDate:    startDt2,
		EndDate:      endDt2,
		CreatorID:    1,
		NotifyBefore: time.Hour * 48,
	})

	r := httptest.NewRequest(
		"POST",
		"http://localhost:8080/event/listByNotifyDate?notify_date=2025-06-13",
		nil,
	)

	w := httptest.NewRecorder()
	server.GetListByNotifyDateHandler(w, r)

	resp := w.Result()
	resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	buf := new(strings.Builder)
	_, err = io.Copy(buf, resp.Body)

	require.Nil(t, err)
	require.Equal(
		t,
		//nolint: all
		"[{\"id\":\"2\",\"title\":\"Test 2\",\"description\":\"\",\"start_dt\":\"2025-06-12T00:00:00Z\",\"end_dt\":\"2025-06-15T00:00:00Z\",\"creator_id\":1,\"notify_before\":172800000000000}]",
		buf.String(),
	)
}
