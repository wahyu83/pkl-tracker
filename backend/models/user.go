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
	Role         string     `gorm:"type:varchar(20);not null;default:'student';check:role IN ('student','teacher','dudi','admin','admin_jurusan')" json:"role"`
	NisNipNik    string     `gorm:"size:50;uniqueIndex;not null" json:"nis_nip_nik"`
	Jurusan      string     `gorm:"size:100" json:"jurusan,omitempty"`
	DudiID       *uuid.UUID `gorm:"type:uuid" json:"dudi_id,omitempty"`
	DUDI         *DUDI      `gorm:"foreignKey:DudiID" json:"dudi,omitempty"`
	TeacherID    *uuid.UUID `gorm:"type:uuid" json:"teacher_id,omitempty"`
	Teacher      *User      `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
