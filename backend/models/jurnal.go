package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Jurnal struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID        uuid.UUID  `gorm:"type:uuid;not null;index" json:"student_id"`
	Student          *User      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Date             time.Time  `gorm:"type:date;not null" json:"date"`
	Activity         string     `gorm:"type:text;not null" json:"activity"`
	DocumentationURL string     `gorm:"type:text" json:"documentation_url"`
	Reflection       string     `gorm:"type:text" json:"reflection"`
	TeacherComment   string     `gorm:"type:text" json:"teacher_comment"`
	DudiComment      string     `gorm:"type:text" json:"dudi_comment"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

func (j *Jurnal) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}
