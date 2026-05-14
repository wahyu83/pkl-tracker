package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Absensi struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID uuid.UUID `gorm:"type:uuid;not null;index" json:"student_id"`
	Student   *User     `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	Latitude  float64   `gorm:"type:decimal(10,8)" json:"latitude"`
	Longitude float64   `gorm:"type:decimal(11,8)" json:"longitude"`
	PhotoURL  string    `gorm:"type:text" json:"photo_url"`
	Status    string    `gorm:"type:varchar(20);not null;default:'hadir';check:status IN ('hadir','terlambat','izin','sakit')" json:"status"`
	IsVerified bool     `gorm:"default:false" json:"is_verified"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *Absensi) BeforeCreate(tx *gorm.DB) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return nil
}
