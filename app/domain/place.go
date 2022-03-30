package domain

import "time"

type Place struct {
	ID        uint       `json:"id"`
	Place     string     `json:"place"`
	Date      time.Time  `json:"date"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
