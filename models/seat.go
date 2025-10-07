package models

import "time"

type Seat struct {
	ID          int        `json:"id"`
	ScheduleID  int        `json:"schedule_id"`
	SeatNumber  string     `json:"seat_number"`
	Status      string     `json:"status"` // available, locked, sold
	LockedUntil *time.Time `json:"locked_until,omitempty"`
	SoldAt      *time.Time `json:"sold_at,omitempty"`
}
