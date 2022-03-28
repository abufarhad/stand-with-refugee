package domain

import "time"

type Specialization struct {
	ID             uint       `json:"id"`
	Specialization string     `json:"specialization"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}
