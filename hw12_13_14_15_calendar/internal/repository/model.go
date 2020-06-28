package repository

import "time"

type EventID uint64

type Event struct {
	ID          EventID       `json:"id"`
	Title       string        `json:"title"`
	Datetime    time.Time     `json:"datetime"`
	Duration    time.Duration `json:"duration"`
	Description string        `json:"description"`
	OwnerID     uint64        `json:"owner_id"`
}

type Notice struct {
	ID       EventID   `json:"id"`
	Title    string    `json:"title"`
	Datetime time.Time `json:"datetime"`
	OwnerID  uint64    `json:"owner_id"`
}
