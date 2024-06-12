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
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartDate    time.Time     `json:"startDt"`
	EndDate      time.Time     `json:"endDt"`
	CreatorID    int           `json:"creatorId"`
	NotifyBefore time.Duration `json:"notifyBefore"`
	Notified     bool
}
