package internalhttp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
)

var (
	ErrServerStarted    = errors.New("server already started")
	ErrServerNotStarted = errors.New("server not started")
)

type Server struct {
	host    string
	port    string
	timeout time.Duration
	server  *http.Server
	mux     *http.ServeMux
	app     Application
	logger  Logger
}

type Logger interface {
	Debug(msg string, params ...any)
	Info(msg string)
	Error(msg string)
}

type Application interface {
	CreateEvent(
		ctx context.Context,
		id, title string,
		startDate time.Time,
		endDate time.Time,
		notifyBefore time.Duration,
	) error

	UpdateEvent(
		ctx context.Context,
		eventID string,
		title string,
		startDate time.Time,
		endDate time.Time,
		notifyBefore time.Duration,
	) error
	DeleteEvent(ctx context.Context, eventID string) error
	GetEvent(ctx context.Context, id string) (storage.Event, error)
	GetEventsListByDates(ctx context.Context, from *time.Time, to *time.Time) []storage.Event
	GetEventsForNotify(ctx context.Context, notifyDate string) []storage.Event
	GetEventsOnDate(ctx context.Context, date time.Time) []storage.Event
	GetEventsOnWeek(ctx context.Context, weekStartDate time.Time) []storage.Event
	GetEventsOnMonth(ctx context.Context, monthStartDate time.Time) []storage.Event
}

func NewServer(logg Logger, app Application, host string, port string, timeout time.Duration) *Server {
	server := &Server{
		host:    host,
		port:    port,
		timeout: timeout,
		mux:     http.NewServeMux(),
		app:     app,
		logger:  logg,
	}

	server.AddRoute("/isready", server.IsReady)
	server.AddRoute("/event/create", server.CreateEventHandler)
	server.AddRoute("/event/update", server.UpdateEventHandler)
	server.AddRoute("/event/delete", server.DeleteEventHandler)
	server.AddRoute("/event/get", server.GetEventHandler)
	server.AddRoute("/event/listByDates", server.GetListByDatesHandler)
	server.AddRoute("/event/listByNotifyDate", server.GetListByNotifyDateHandler)
	server.AddRoute("/event/listOnDate", server.GetListOnDateHandler)
	server.AddRoute("/event/listOnWeek", server.GetListOnWeekHandler)
	server.AddRoute("/event/listOnMonth", server.GetListOnMonthHandler)

	return server
}

func (s *Server) Start(_ context.Context) error {
	if s.server != nil {
		return ErrServerStarted
	}
	address := net.JoinHostPort(s.host, s.port)
	s.server = &http.Server{
		Addr:              address,
		ReadHeaderTimeout: s.timeout,
		ReadTimeout:       s.timeout,
		Handler:           loggingMiddleware(s.mux),
	}

	err := s.server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	if s.server == nil {
		return ErrServerNotStarted
	}

	return s.server.Close()
}

func (s *Server) AddRoute(route string, handlerFunc http.HandlerFunc) {
	s.mux.HandleFunc(route, handlerFunc)
}

func (s *Server) IsReady(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		s.logger.Error(err.Error())
		s.internalError(w, err)
		return
	}
}

func (s *Server) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")

	startDt, err := time.Parse(time.DateOnly, r.FormValue("start_dt"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("start_dt: "+err.Error()))
		return
	}

	endDt, err := time.Parse(time.DateOnly, r.FormValue("end_dt"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("end_dt: "+err.Error()))
		return
	}

	notifyBefore, err := time.ParseDuration(r.FormValue("notify_before"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("notify_before: "+err.Error()))
		return
	}

	err = s.app.CreateEvent(
		r.Context(),
		id,
		title,
		startDt,
		endDt,
		notifyBefore,
	)

	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, err)
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")
	title := r.PostFormValue("title")

	startDt, err := time.Parse(time.DateOnly, r.PostFormValue("start_dt"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("start_dt: "+err.Error()))
	}

	endDt, err := time.Parse(time.DateOnly, r.PostFormValue("end_dt"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("end_dt: "+err.Error()))
		return
	}

	notifyBefore, err := time.ParseDuration(r.PostFormValue("notify_before"))
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, errors.New("notify_before: "+err.Error()))
		return
	}

	err = s.app.UpdateEvent(
		r.Context(),
		id,
		title,
		startDt,
		endDt,
		notifyBefore,
	)

	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, err)
	}
}

func (s *Server) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PostFormValue("id")

	err := s.app.DeleteEvent(r.Context(), id)
	if err != nil {
		s.logger.Error(err.Error())
		s.badRequest(w, err)
	}
}

func (s *Server) GetEventHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("id") {
		s.badRequest(w, errors.New("id not passed"))
		return
	}

	id := r.URL.Query().Get("id")

	event, err := s.app.GetEvent(
		r.Context(),
		id,
	)
	if err != nil {
		s.logger.Error(err.Error())
		s.internalError(w, err)
	}

	json, err := event.MarshalJSON()
	if err != nil {
		s.logger.Error(err.Error())
		s.internalError(w, err)
	}

	_, writeErr := w.Write(json)
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) GetListByDatesHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("from") {
		s.badRequest(w, errors.New("from not passed"))
		return
	}

	if !r.URL.Query().Has("to") {
		s.badRequest(w, errors.New("to not passed"))
		return
	}

	fromDt, err := time.Parse(time.DateOnly, r.URL.Query().Get("from"))
	if err != nil {
		s.badRequest(w, errors.New("from invalid format"))
	}

	toDt, err := time.Parse(time.DateOnly, r.URL.Query().Get("to"))
	if err != nil {
		s.badRequest(w, errors.New("to invalid format"))
	}

	events := s.app.GetEventsListByDates(
		r.Context(),
		&fromDt,
		&toDt,
	)

	b := strings.Builder{}
	_, err = b.WriteString("[")
	if err != nil {
		s.internalError(w, err)
	}

	for i := range events {
		jsonEvent, err := events[i].MarshalJSON()
		if err != nil {
			s.internalError(w, err)
		}

		_, err = b.WriteString(string(jsonEvent))
		if err != nil {
			s.internalError(w, err)
		}
	}

	_, err = b.WriteString("]")
	if err != nil {
		s.internalError(w, err)
	}

	_, writeErr := w.Write([]byte(b.String()))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) GetListByNotifyDateHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("notify_date") {
		s.badRequest(w, errors.New("notify_date not passed"))
		return
	}

	events := s.app.GetEventsForNotify(r.Context(), r.URL.Query().Get("notify_date"))

	b := strings.Builder{}
	_, err := b.WriteString("[")
	if err != nil {
		s.internalError(w, err)
	}

	for i := range events {
		jsonEvent, err := events[i].MarshalJSON()
		if err != nil {
			s.internalError(w, err)
		}

		_, err = b.WriteString(string(jsonEvent))
		if err != nil {
			s.internalError(w, err)
		}
	}

	_, err = b.WriteString("]")
	if err != nil {
		s.internalError(w, err)
	}

	_, writeErr := w.Write([]byte(b.String()))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) GetListOnDateHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("date") {
		s.badRequest(w, errors.New("date not passed"))
		return
	}

	date, err := time.Parse(time.DateOnly, r.URL.Query().Get("date"))
	if err != nil {
		s.badRequest(w, errors.New("date invalid format"))
	}

	events := s.app.GetEventsOnDate(r.Context(), date)

	jsonStr, err := buildEventsJSON(events)
	if err != nil {
		s.internalError(w, err)
		return
	}

	_, writeErr := w.Write([]byte(jsonStr))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) GetListOnWeekHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("weekStartDate") {
		s.badRequest(w, errors.New("weekStartDate not passed"))
		return
	}

	weekStartDate, err := time.Parse(time.DateOnly, r.URL.Query().Get("weekStartDate"))
	if err != nil {
		s.badRequest(w, errors.New("weekStartDate invalid format"))
	}

	events := s.app.GetEventsOnWeek(r.Context(), weekStartDate)

	jsonStr, err := buildEventsJSON(events)
	if err != nil {
		s.internalError(w, err)
		return
	}

	_, writeErr := w.Write([]byte(jsonStr))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) GetListOnMonthHandler(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("monthStartDate") {
		s.badRequest(w, errors.New("monthStartDate not passed"))
		return
	}

	monthStartDate, err := time.Parse(time.DateOnly, r.URL.Query().Get("monthStartDate"))
	if err != nil {
		s.badRequest(w, errors.New("monthStartDate invalid format"))
	}

	events := s.app.GetEventsOnMonth(r.Context(), monthStartDate)

	jsonStr, err := buildEventsJSON(events)
	if err != nil {
		s.internalError(w, err)
		return
	}

	_, writeErr := w.Write([]byte(jsonStr))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func buildEventsJSON(events []storage.Event) (string, error) {
	b := strings.Builder{}
	_, err := b.WriteString("[")
	if err != nil {
		return "", err
	}

	for i := range events {
		jsonEvent, err := events[i].MarshalJSON()
		if err != nil {
			return "", err
		}

		_, err = b.WriteString(string(jsonEvent))
		if err != nil {
			return "", err
		}
	}

	_, err = b.WriteString("]")
	if err != nil {
		return "", err
	}

	return b.String(), nil
}

func (s *Server) internalError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	_, writeErr := w.Write([]byte(err.Error()))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}

func (s *Server) badRequest(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusBadRequest)
	_, writeErr := w.Write([]byte(err.Error()))
	if writeErr != nil {
		s.logger.Error(writeErr.Error())
	}
}
