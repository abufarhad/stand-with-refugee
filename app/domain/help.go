package domain

import (
	"gorm.io/datatypes"
	"time"
)

type Help struct {
	ID                     uint           `json:"id"`
	Name                   string         `json:"name"`
	Phone                  string         `json:"phone"`
	Gender                 string         `json:"gender"`
	Age                    uint           `json:"age"`
	PlaceID                uint           `json:"place_id"`
	Place                  Place          `gorm:"foreignKey:PlaceID" json:"-"`
	Symptoms               datatypes.JSON `json:"symptoms"`
	SpecializationTypeNeed string         `json:"specialization_type_need"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              *time.Time     `json:"deleted_at"`
}
