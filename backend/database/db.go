package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"pkl-tracker/config"
	"pkl-tracker/models"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.DUDI{},
		&models.Absensi{},
		&models.Jurnal{},
		&models.Penilaian{},
		&models.Periode{},
		&models.Jurusan{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	migrateForeignKeyCascade()

	log.Println("Database connected and migrated successfully")
}

func migrateForeignKeyCascade() {
	cascades := []struct {
		table    string
		constraint string
		column     string
	}{
		{"absensis", "fk_absensis_student", "student_id"},
		{"jurnals", "fk_jurnals_student", "student_id"},
		{"penilaians", "fk_penilaians_student", "student_id"},
	}

	for _, fk := range cascades {
		var exists bool
		DB.Raw(`SELECT EXISTS (
			SELECT 1 FROM information_schema.table_constraints
			WHERE constraint_name = ? AND table_name = ?
		)`, fk.constraint, fk.table).Scan(&exists)

		if !exists {
			continue
		}

		DB.Exec("ALTER TABLE " + fk.table + " DROP CONSTRAINT " + fk.constraint)

		DB.Exec("ALTER TABLE " + fk.table +
			" ADD CONSTRAINT " + fk.constraint +
			" FOREIGN KEY (" + fk.column + ") REFERENCES users(id) ON DELETE CASCADE")
	}
}
