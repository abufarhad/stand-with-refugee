package domain

import "time"

type Help struct {
	ID        uint       `json:"id"`
	Name      string     `json:"name"`
	Phone     string     `json:"phone"`
	Gender    string     `json:"gender"`
	Age       uint       `json:"age"`
	Place     string     `json:"place"`
	Symptoms  []string   `json:"symptoms"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
