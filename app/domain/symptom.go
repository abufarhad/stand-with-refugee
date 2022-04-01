package domain

import "time"

type Symptom struct {
	ID        uint       `json:"id"`
	Symptom   string     `json:"symptom"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
