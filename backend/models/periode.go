package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Periode struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TahunPelajaran string    `gorm:"size:20;not null" json:"tahun_pelajaran"`
	Semester       string    `gorm:"type:varchar(10);not null;check:semester IN ('ganjil','genap')" json:"semester"`
	IsActive       bool      `gorm:"default:false" json:"is_active"`
	StartDate      time.Time `gorm:"type:date" json:"start_date"`
	EndDate        time.Time `gorm:"type:date" json:"end_date"`
	CreatedAt      time.Time `json:"created_at"`
}

func (p *Periode) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return nil
}
