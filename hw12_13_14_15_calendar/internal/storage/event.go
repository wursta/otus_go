package storage

import (
	"errors"
	"time"
)

var (
	ErrReadEventNotExists     = errors.New("read event: passed ID not exists")
	ErrCreateEventIDExists    = errors.New("create event: passed ID exists")
	ErrUpdateEventIDNotExists = errors.New("update event: event with passed ID not exists")
)

type Event struct {
	ID           string
	Title        string
	Description  string
	StartDate    time.Time
	EndDate      time.Time
	CreatorID    int
	NotifyBefore time.Duration
}
