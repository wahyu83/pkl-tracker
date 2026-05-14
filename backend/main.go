package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pkl-tracker/config"
	"pkl-tracker/database"
	"pkl-tracker/handlers"
	"pkl-tracker/middleware"
	"pkl-tracker/models"
)

func main() {
	cfg := config.Load()
	middleware.Init(cfg)

	database.Connect(cfg)
	seedDatabase()

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	api := r.Group("/api")

	authHandler := handlers.NewAuthHandler(cfg)

	api.POST("/login", authHandler.Login)
	api.POST("/register", authHandler.Register)

	protected := api.Group("")
	protected.Use(middleware.AuthRequired())
	{
		protected.GET("/me", authHandler.Me)
		protected.POST("/change-password", authHandler.ChangePassword)

		absensiHandler := handlers.NewAbsensiHandler()
		protected.POST("/absensi", absensiHandler.Create)
		protected.GET("/absensi/history", absensiHandler.History)
		protected.PUT("/absensi/:id/verify", absensiHandler.Verify)

		jurnalHandler := handlers.NewJurnalHandler()
		protected.POST("/jurnal", jurnalHandler.Create)
		protected.GET("/jurnal", jurnalHandler.List)
		protected.GET("/jurnal/:id", jurnalHandler.GetByID)
		protected.PUT("/jurnal/:id", jurnalHandler.Update)
		protected.POST("/jurnal/comment", jurnalHandler.Comment)

		penilaianHandler := handlers.NewPenilaianHandler()
		protected.POST("/nilai", penilaianHandler.CreateOrUpdate)
		protected.GET("/nilai", penilaianHandler.List)
		protected.GET("/nilai/:studentId", penilaianHandler.GetByStudent)

		reportHandler := handlers.NewReportHandler()
		protected.GET("/report/absensi", reportHandler.AbsensiReport)
		protected.GET("/report/jurnal", reportHandler.JurnalReport)
		protected.GET("/report/nilai", reportHandler.NilaiReport)

		importHandler := handlers.NewImportHandler()
		protected.POST("/import/siswa", importHandler.ImportSiswa)
		protected.POST("/import/guru", importHandler.ImportGuru)
		protected.POST("/import/dudi", importHandler.ImportDudi)
		protected.POST("/import/instruktur-dudi", importHandler.ImportInstrukturDudi)

		exportHandler := handlers.NewExportHandler()
		protected.GET("/export/absensi", exportHandler.ExportAbsensi)
		protected.GET("/export/jurnal", exportHandler.ExportJurnal)
		protected.GET("/export/nilai", exportHandler.ExportNilai)
	}

	log.Printf("Server starting on port %s", cfg.ServerPort)
	r.Run(":" + cfg.ServerPort)
}

func seedDatabase() {
	var count int64
	database.DB.Model(&models.User{}).Count(&count)
	if count > 0 {
		return
	}

	dudiID1 := uuid.New()
	dudiID2 := uuid.New()

	dudi1 := models.DUDI{
		ID: dudiID1, CompanyName: "PT. Teknologi Maju", Address: "Jl. Sudirman No. 123, Jakarta Pusat",
		Latitude: -6.2088, Longitude: 106.8456, RadiusAllowed: 500, PicName: "Hendra Gunawan", Phone: "021-5551234",
	}
	dudi2 := models.DUDI{
		ID: dudiID2, CompanyName: "PT. Sejahtera Abadi", Address: "Jl. Gatot Subroto No. 45, Jakarta Selatan",
		Latitude: -6.2297, Longitude: 106.8243, RadiusAllowed: 300, PicName: "Ratna Dewi", Phone: "021-5555678",
	}
	database.DB.Create(&dudi1)
	database.DB.Create(&dudi2)

	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	teacherPass, _ := bcrypt.GenerateFromPassword([]byte("guru123"), bcrypt.DefaultCost)
	studentPass, _ := bcrypt.GenerateFromPassword([]byte("siswa123"), bcrypt.DefaultCost)
	dudiPass, _ := bcrypt.GenerateFromPassword([]byte("dudi123"), bcrypt.DefaultCost)

	adminID := uuid.New()
	teacherID := uuid.New()
	student1ID := uuid.New()
	student2ID := uuid.New()
	dudiUserID := uuid.New()

	users := []models.User{
		{ID: adminID, FullName: "Admin Utama", Email: "admin@pkl.local", PasswordHash: string(adminPass), Role: "admin", NisNipNik: "ADM-001"},
		{ID: teacherID, FullName: "Budi Santoso, S.Kom", Email: "budi@pkl.local", PasswordHash: string(teacherPass), Role: "teacher", NisNipNik: "19850101"},
		{ID: student1ID, FullName: "Ahmad Rizky", Email: "ahmad@pkl.local", PasswordHash: string(studentPass), Role: "student", NisNipNik: "20230001", DudiID: &dudiID1},
		{ID: student2ID, FullName: "Siti Nurhaliza", Email: "siti@pkl.local", PasswordHash: string(studentPass), Role: "student", NisNipNik: "20230002", DudiID: &dudiID2},
		{ID: dudiUserID, FullName: "PT. Teknologi Maju", Email: "info@teknologimaju.id", PasswordHash: string(dudiPass), Role: "dudi", NisNipNik: "D-001", DudiID: &dudiID1},
	}

	for _, u := range users {
		database.DB.Create(&u)
	}

	now := time.Now()
	absensiList := []models.Absensi{
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day(), 7, 45, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Status: "hadir", IsVerified: true},
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-1, 8, 15, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Status: "terlambat", IsVerified: true},
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-2, 7, 30, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Status: "hadir", IsVerified: true},
		{StudentID: student2ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-1, 7, 50, 0, 0, time.UTC), Latitude: -6.2297, Longitude: 106.8243, Status: "hadir", IsVerified: true},
	}

	for _, a := range absensiList {
		database.DB.Create(&a)
	}

	jurnalList := []models.Jurnal{
		{StudentID: student1ID, Date: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
			Activity: "Mempelajari framework Laravel dan membuat CRUD untuk modul inventaris.", Reflection: "Belajar banyak tentang MVC pattern."},
		{StudentID: student1ID, Date: time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			Activity: "Debugging aplikasi internal, memperbaiki bug modul pelaporan.",
			TeacherComment: "Terus tingkatkan!"},
		{StudentID: student2ID, Date: time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			Activity: "Membantu tim network maintenance server.",
			DudiComment: "Siswa menunjukkan antusiasme baik."},
	}

	for _, j := range jurnalList {
		database.DB.Create(&j)
	}

	log.Println("Seed data created successfully!")
	log.Println("Test accounts:")
	log.Println("  Admin:  NIS=ADM-001  Password=admin123")
	log.Println("  Guru:   NIP=19850101  Password=guru123")
	log.Println("  Siswa:  NIS=20230001  Password=siswa123")
	log.Println("  DUDI:   NIK=D-001     Password=dudi123")
}
