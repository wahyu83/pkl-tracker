package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	FullName     string     `gorm:"size:255;not null" json:"full_name"`
	Email        string     `gorm:"size:255;uniqueIndex;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         string     `gorm:"type:varchar(20);not null;default:'student';check:role IN ('student','teacher','dudi','admin')" json:"role"`
	NisNipNik    string     `gorm:"size:50;uniqueIndex;not null" json:"nis_nip_nik"`
	DudiID       *uuid.UUID `gorm:"type:uuid" json:"dudi_id,omitempty"`
	DUDI         *DUDI      `gorm:"foreignKey:DudiID" json:"dudi,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
