package models

import "time"

type Schedule struct {
	ID           int       `json:"id"`
	MovieTitle   string    `json:"movie_title"`
	CinemaBranch string    `json:"cinema_branch"`
	City         string    `json:"city"`
	ShowTime     time.Time `json:"show_time"`
	TotalSeats   int       `json:"total_seats"`
	Status       string    `json:"status"`
}
