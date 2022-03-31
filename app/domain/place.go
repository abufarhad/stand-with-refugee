package domain

import "time"

type Place struct {
	ID        uint       `json:"id"`
	PlaceName string     `json:"place_name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
