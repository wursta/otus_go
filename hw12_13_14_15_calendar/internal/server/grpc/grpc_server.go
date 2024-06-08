package internalgrpc

import (
	"context"
	"net"
	"time"

	calendarpb "github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/wursta/otus_go/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

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

type Server struct {
	logger Logger
	app    Application
	port   string
	server *grpc.Server
	calendarpb.UnimplementedCalendarServer
}

func NewServer(logg Logger, app Application, port string) *Server {
	return &Server{
		logger: logg,
		app:    app,
		port:   port,
	}
}

func (s *Server) Start(_ context.Context) error {
	lsn, err := net.Listen("tcp4", net.JoinHostPort("localhost", s.port))
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}

	s.server = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			loggingMiddleware(),
		),
	)

	calendarpb.RegisterCalendarServer(s.server, s)
	reflection.Register(s.server)

	if err := s.server.Serve(lsn); err != nil {
		s.logger.Error(err.Error())
		return err
	}

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	s.server.GracefulStop()
	return nil
}

func (s *Server) Create(ctx context.Context, r *calendarpb.CreateRequest) (*calendarpb.CreateResult, error) {
	id := r.GetId()
	title := r.GetTitle()
	startDt := r.GetStartDt().AsTime()
	endDt := r.GetEndDt().AsTime()
	notifyBefore := r.GetNotifyBefore().AsDuration()

	err := s.app.CreateEvent(
		ctx,
		id,
		title,
		startDt,
		endDt,
		notifyBefore,
	)
	if err != nil {
		return nil, err
	}

	return &calendarpb.CreateResult{}, nil
}

func (s *Server) Update(ctx context.Context, r *calendarpb.UpdateRequest) (*calendarpb.UpdateResult, error) {
	id := r.GetEventId()
	title := r.GetTitle()
	startDt := r.GetStartDt().AsTime()
	endDt := r.GetEndDt().AsTime()
	notifyBefore := r.GetNotifyBefore().AsDuration()

	err := s.app.UpdateEvent(
		ctx,
		id,
		title,
		startDt,
		endDt,
		notifyBefore,
	)
	if err != nil {
		return nil, err
	}

	return &calendarpb.UpdateResult{}, nil
}

func (s *Server) Delete(ctx context.Context, r *calendarpb.DeleteRequest) (*calendarpb.DeleteResult, error) {
	eventID := r.GetEventId()

	err := s.app.DeleteEvent(ctx, eventID)
	if err != nil {
		return nil, err
	}

	return &calendarpb.DeleteResult{}, nil
}

func (s *Server) Get(ctx context.Context, r *calendarpb.GetRequest) (*calendarpb.GetResult, error) {
	id := r.GetId()

	event, err := s.app.GetEvent(ctx, id)
	if err != nil {
		return nil, err
	}

	return &calendarpb.GetResult{
		Id:           event.ID,
		Title:        event.Title,
		StartDt:      timestamppb.New(event.StartDate),
		EndDt:        timestamppb.New(event.EndDate),
		NotifyBefore: durationpb.New(event.NotifyBefore),
	}, nil
}

func (s *Server) GetEventsListByDates(
	ctx context.Context,
	r *calendarpb.GetEventsListByDatesRequest,
) (*calendarpb.GetEventsListByDatesResult, error) {
	from := r.GetFrom().AsTime()
	to := r.GetTo().AsTime()

	events := s.app.GetEventsListByDates(ctx, &from, &to)

	resultsList := []*calendarpb.GetResult{}

	for i := range events {
		resultsList = append(resultsList, &calendarpb.GetResult{
			Id:           events[i].ID,
			Title:        events[i].Title,
			StartDt:      timestamppb.New(events[i].StartDate),
			EndDt:        timestamppb.New(events[i].EndDate),
			NotifyBefore: durationpb.New(events[i].NotifyBefore),
		})
	}

	return &calendarpb.GetEventsListByDatesResult{
		List: resultsList,
	}, nil
}

func (s *Server) GetEventsForNotify(
	ctx context.Context,
	r *calendarpb.GetEventsForNotifyRequest,
) (*calendarpb.GetEventsForNotifyResult, error) {
	notifyDate := r.GetNotifyDate()
	events := s.app.GetEventsForNotify(ctx, notifyDate)

	resultsList := []*calendarpb.GetResult{}

	for i := range events {
		resultsList = append(resultsList, &calendarpb.GetResult{
			Id:           events[i].ID,
			Title:        events[i].Title,
			StartDt:      timestamppb.New(events[i].StartDate),
			EndDt:        timestamppb.New(events[i].EndDate),
			NotifyBefore: durationpb.New(events[i].NotifyBefore),
		})
	}

	return &calendarpb.GetEventsForNotifyResult{
		List: resultsList,
	}, nil
}

func (s *Server) GetEventsListOnDate(
	ctx context.Context,
	r *calendarpb.GetEventsListOnDateRequest,
) (*calendarpb.GetEventsListOnDateResult, error) {
	dayDate := r.GetDayDate().AsTime()

	events := s.app.GetEventsOnDate(ctx, dayDate)

	resultsList := []*calendarpb.GetResult{}

	for i := range events {
		resultsList = append(resultsList, &calendarpb.GetResult{
			Id:           events[i].ID,
			Title:        events[i].Title,
			StartDt:      timestamppb.New(events[i].StartDate),
			EndDt:        timestamppb.New(events[i].EndDate),
			NotifyBefore: durationpb.New(events[i].NotifyBefore),
		})
	}

	return &calendarpb.GetEventsListOnDateResult{
		List: resultsList,
	}, nil
}

func (s *Server) GetEventsListOnWeek(
	ctx context.Context,
	r *calendarpb.GetEventsListOnWeekRequest,
) (*calendarpb.GetEventsListOnWeekResult, error) {
	weekStartDate := r.GetWeekStartDate().AsTime()

	events := s.app.GetEventsOnWeek(ctx, weekStartDate)

	resultsList := []*calendarpb.GetResult{}

	for i := range events {
		resultsList = append(resultsList, &calendarpb.GetResult{
			Id:           events[i].ID,
			Title:        events[i].Title,
			StartDt:      timestamppb.New(events[i].StartDate),
			EndDt:        timestamppb.New(events[i].EndDate),
			NotifyBefore: durationpb.New(events[i].NotifyBefore),
		})
	}

	return &calendarpb.GetEventsListOnWeekResult{
		List: resultsList,
	}, nil
}

func (s *Server) GetEventsListOnMonth(
	ctx context.Context,
	r *calendarpb.GetEventsListOnMonthRequest,
) (*calendarpb.GetEventsListOnMonthResult, error) {
	weekStartDate := r.GetMonthStartDate().AsTime()

	events := s.app.GetEventsOnMonth(ctx, weekStartDate)

	resultsList := []*calendarpb.GetResult{}

	for i := range events {
		resultsList = append(resultsList, &calendarpb.GetResult{
			Id:           events[i].ID,
			Title:        events[i].Title,
			StartDt:      timestamppb.New(events[i].StartDate),
			EndDt:        timestamppb.New(events[i].EndDate),
			NotifyBefore: durationpb.New(events[i].NotifyBefore),
		})
	}

	return &calendarpb.GetEventsListOnMonthResult{
		List: resultsList,
	}, nil
}
