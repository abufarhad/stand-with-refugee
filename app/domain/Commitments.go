package domain

import "time"

type Commitments struct {
	ID       uint `json:"id"`
	DoctorId uint `json:"doctor_id"`
	PlaceID  uint `json:"place_id"`
	//Place     Place      `gorm:"foreignKey:PlaceID" json:"-"`
	Point     uint       `json:"point"`
	Date      time.Time  `json:"date"`
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
