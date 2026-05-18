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

	r := gin.Default()

	r.Use(spaMiddleware())

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

	r.Static("/uploads", "./uploads")

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
		protected.GET("/absensi/status", absensiHandler.Status)
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

		adminHandler := handlers.NewAdminHandler()
		protected.GET("/admin/users", adminHandler.ListUsers)
		protected.POST("/admin/users", adminHandler.CreateUser)
		protected.PUT("/admin/users/:id", adminHandler.UpdateUser)
		protected.DELETE("/admin/users/:id", adminHandler.DeleteUser)
		protected.POST("/admin/users/bulk-delete", adminHandler.BulkDeleteUsers)
		protected.GET("/admin/dudi", adminHandler.ListDUDI)
		protected.POST("/admin/dudi", adminHandler.CreateDUDI)
		protected.PUT("/admin/dudi/:id", adminHandler.UpdateDUDI)
		protected.DELETE("/admin/dudi/:id", adminHandler.DeleteDUDI)
		protected.GET("/admin/periode", adminHandler.ListPeriode)
		protected.GET("/admin/periode/active", adminHandler.ActivePeriode)
		protected.POST("/admin/periode", adminHandler.CreatePeriode)
		protected.PUT("/admin/periode/:id", adminHandler.UpdatePeriode)
		protected.PUT("/admin/periode/:id/activate", adminHandler.ActivatePeriode)
		protected.DELETE("/admin/periode/:id", adminHandler.DeletePeriode)
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
		Latitude: -6.2088, Longitude: 106.8456, RadiusAllowed: 500, PicName: "Hendra Gunawan", Phone: "021-5551234", Jurusan: "RPL",
	}
	dudi2 := models.DUDI{
		ID: dudiID2, CompanyName: "PT. Sejahtera Abadi", Address: "Jl. Gatot Subroto No. 45, Jakarta Selatan",
		Latitude: -6.2297, Longitude: 106.8243, RadiusAllowed: 300, PicName: "Ratna Dewi", Phone: "021-5555678", Jurusan: "TKJ",
	}
	database.DB.Create(&dudi1)
	database.DB.Create(&dudi2)

	adminPass, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	teacherPass, _ := bcrypt.GenerateFromPassword([]byte("guru123"), bcrypt.DefaultCost)
	studentPass, _ := bcrypt.GenerateFromPassword([]byte("siswa123"), bcrypt.DefaultCost)
	dudiPass, _ := bcrypt.GenerateFromPassword([]byte("dudi123"), bcrypt.DefaultCost)
	adminJurusanPass, _ := bcrypt.GenerateFromPassword([]byte("jurusan123"), bcrypt.DefaultCost)

	adminID := uuid.New()
	teacherID := uuid.New()
	student1ID := uuid.New()
	student2ID := uuid.New()
	dudiUserID := uuid.New()
	adminJurusanID := uuid.New()

	users := []models.User{
		{ID: adminID, FullName: "Admin Utama", Email: "admin@pkl.local", PasswordHash: string(adminPass), Role: "admin", NisNipNik: "ADM-001"},
		{ID: teacherID, FullName: "Budi Santoso, S.Kom", Email: "budi@pkl.local", PasswordHash: string(teacherPass), Role: "teacher", NisNipNik: "19850101", Jurusan: "RPL"},
		{ID: student1ID, FullName: "Ahmad Rizky", Email: "ahmad@pkl.local", PasswordHash: string(studentPass), Role: "student", NisNipNik: "20230001", DudiID: &dudiID1, Jurusan: "RPL"},
		{ID: student2ID, FullName: "Siti Nurhaliza", Email: "siti@pkl.local", PasswordHash: string(studentPass), Role: "student", NisNipNik: "20230002", DudiID: &dudiID2, Jurusan: "TKJ"},
		{ID: dudiUserID, FullName: "PT. Teknologi Maju", Email: "info@teknologimaju.id", PasswordHash: string(dudiPass), Role: "dudi", NisNipNik: "D-001", DudiID: &dudiID1, Jurusan: "RPL"},
		{ID: adminJurusanID, FullName: "Admin RPL", Email: "rpl@pkl.local", PasswordHash: string(adminJurusanPass), Role: "admin_jurusan", NisNipNik: "ADM-RPL", Jurusan: "RPL"},
	}

	for _, u := range users {
		database.DB.Create(&u)
	}

	now := time.Now()
	absensiList := []models.Absensi{
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day(), 7, 45, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Type: "masuk", Status: "hadir", IsVerified: true},
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-1, 8, 15, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Type: "masuk", Status: "hadir", IsVerified: true},
		{StudentID: student1ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-2, 7, 30, 0, 0, time.UTC), Latitude: -6.2088, Longitude: 106.8456, Type: "masuk", Status: "hadir", IsVerified: true},
		{StudentID: student2ID, Timestamp: time.Date(now.Year(), now.Month(), now.Day()-1, 7, 50, 0, 0, time.UTC), Latitude: -6.2297, Longitude: 106.8243, Type: "masuk", Status: "hadir", IsVerified: true},
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
	log.Println("  Admin:          NIS=ADM-001  Password=admin123")
	log.Println("  Admin Jurusan:  NIS=ADM-RPL  Password=jurusan123")
	log.Println("  Guru:           NIP=19850101  Password=guru123")
	log.Println("  Siswa:          NIS=20230001  Password=siswa123")
	log.Println("  DUDI:           NIK=D-001     Password=dudi123")
}

func seedPeriode() {
	var count int64
	database.DB.Model(&models.Periode{}).Count(&count)
	if count > 0 {
		return
	}

	periods := []models.Periode{
		{
			TahunPelajaran: "2025/2026",
			Semester:       "ganjil",
			IsActive:       true,
			StartDate:      time.Date(2025, 7, 14, 0, 0, 0, 0, time.UTC),
			EndDate:        time.Date(2025, 12, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			TahunPelajaran: "2025/2026",
			Semester:       "genap",
			IsActive:       false,
			StartDate:      time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC),
			EndDate:        time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, p := range periods {
		database.DB.Create(&p)
	}

	log.Println("Period seed data created!")
}
