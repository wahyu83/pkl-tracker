package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Penilaian struct {
	ID                  uuid.UUID  `gorm:"type:uuid;primaryKey" json:"id"`
	StudentID           uuid.UUID  `gorm:"type:uuid;not null;index;constraint:OnDelete:CASCADE" json:"student_id"`
	Student             *User      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	DudiID              uuid.UUID  `gorm:"type:uuid;not null" json:"dudi_id"`
	AttendanceScoreAuto float64    `gorm:"type:decimal(5,2);default:0" json:"attendance_score_auto"`
	Discipline          int        `json:"discipline"`
	Responsibility      int        `json:"responsibility"`
	Teamwork            int        `json:"teamwork"`
	Initiative          int        `json:"initiative"`
	FinalScore          float64    `gorm:"type:decimal(5,2)" json:"final_score"`
	FinalGrade          string     `gorm:"size:2" json:"final_grade"`
	Notes               string     `gorm:"type:text" json:"notes"`
	SubmittedAt         time.Time  `json:"submitted_at"`
}

func (p *Penilaian) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
