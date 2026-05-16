package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DUDI struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	CompanyName   string    `gorm:"size:255;not null" json:"company_name"`
	Address       string    `gorm:"type:text" json:"address"`
	Latitude      float64   `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude     float64   `gorm:"type:decimal(11,8)" json:"longitude"`
	RadiusAllowed int       `gorm:"default:500" json:"radius_allowed"`
	PicName       string    `gorm:"size:255" json:"pic_name"`
	Phone         string    `gorm:"size:20" json:"phone"`
	Jurusan       string    `gorm:"size:100" json:"jurusan,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

func (d *DUDI) BeforeCreate(tx *gorm.DB) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	return nil
}
