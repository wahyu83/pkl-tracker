package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SchoolProfile menyimpan identitas sekolah yang dipakai untuk kop dan
// tanda tangan pada sertifikat PKL. Hanya ada satu baris (singleton).
type SchoolProfile struct {
	ID                uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Name              string    `gorm:"size:255" json:"name"`
	Address           string    `gorm:"type:text" json:"address"`
	City              string    `gorm:"size:100" json:"city"`
	HeadmasterName    string    `gorm:"size:255" json:"headmaster_name"`
	HeadmasterNip     string    `gorm:"size:50" json:"headmaster_nip"`
	CertificatePrefix string    `gorm:"size:100" json:"certificate_prefix"`
	LogoURL           string    `gorm:"type:text" json:"logo_url"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (s *SchoolProfile) BeforeCreate(tx *gorm.DB) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return nil
}
