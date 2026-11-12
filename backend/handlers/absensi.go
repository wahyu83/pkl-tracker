package handlers

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type AbsensiHandler struct{}

func NewAbsensiHandler() *AbsensiHandler {
	return &AbsensiHandler{}
}

type AbsensiRequest struct {
	Type      string  `json:"type" binding:"required"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Status    string  `json:"status"`
	// StudentID diisi saat guru/DUDI melakukan take over absensi siswa.
	StudentID string `json:"student_id"`
}

// resolveStudentForStaff memvalidasi bahwa guru/DUDI berhak mencatat absensi
// untuk siswa tertentu, lalu mengembalikan data siswa tersebut.
// Mengembalikan pesan error bila tidak berhak.
func resolveStudentForStaff(staffID uuid.UUID, role, studentID string) (*models.User, string) {
	studentUID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, "ID siswa tidak valid"
	}

	var student models.User
	if err := database.DB.First(&student, "id = ? AND role = 'student'", studentUID).Error; err != nil {
		return nil, "Siswa tidak ditemukan"
	}

	if role == "teacher" {
		if student.TeacherID == nil || *student.TeacherID != staffID {
			return nil, "Siswa ini tidak dalam bimbingan Anda"
		}
		return &student, ""
	}

	var dudiUser models.User
	if database.DB.First(&dudiUser, "id = ?", staffID).Error != nil || dudiUser.DudiID == nil {
		return nil, "Profil DUDI tidak ditemukan"
	}
	if student.DudiID == nil || *student.DudiID != *dudiUser.DudiID {
		return nil, "Siswa ini tidak berada di DUDI Anda"
	}
	return &student, ""
}

// clientDevice mengembalikan IP dan user agent perangkat yang dipakai absen.
// Untuk entri yang dibuat guru/DUDI (take over / izin / sakit) keduanya
// dikosongkan: perangkat yang dipakai milik guru/DUDI, bukan siswa. Kalau
// disimpan, baris ini akan memicu deteksi "perangkat bersama" palsu pada
// absensi mandiri siswa lain yang lewat jaringan yang sama.
func clientDevice(c *gin.Context, staffEntry bool) (string, string) {
	if staffEntry {
		return "", ""
	}
	return c.ClientIP(), c.Request.UserAgent()
}

func (h *AbsensiHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	uid, _ := uuid.Parse(userID.(string))

	var req AbsensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data absensi tidak valid"})
		return
	}

	if req.Type != "masuk" && req.Type != "pulang" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'masuk' or 'pulang'"})
		return
	}

	// student  : absen untuk dirinya sendiri (butuh lokasi GPS)
	// teacher  : take over absensi siswa bimbingannya (butuh student_id)
	// dudi     : take over absensi siswa di DUDI-nya (butuh student_id)
	var student models.User
	takeOver := false

	switch role {
	case "student":
		if err := database.DB.First(&student, "id = ?", uid).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
			return
		}
		if req.Latitude == 0 && req.Longitude == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude required"})
			return
		}
	case "teacher", "dudi":
		if req.StudentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student_id diperlukan untuk take over absensi"})
			return
		}
		s, msg := resolveStudentForStaff(uid, role.(string), req.StudentID)
		if s == nil {
			status := http.StatusForbidden
			if msg == "Siswa tidak ditemukan" || msg == "ID siswa tidak valid" {
				status = http.StatusBadRequest
			}
			c.JSON(status, gin.H{"error": msg})
			return
		}
		student = *s
		takeOver = true
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	sid := student.ID
	subject := "Anda"
	if takeOver {
		subject = "Siswa"
	}

	// Gunakan waktu server (WIB)
	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	if req.Type == "masuk" {
		// Cek sudah absen masuk hari ini?
		var count int64
		database.DB.Model(&models.Absensi{}).
			Where("student_id = ? AND type = 'masuk' AND timestamp >= ? AND timestamp < ?", sid, todayStart, todayEnd).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": subject + " sudah absen masuk hari ini"})
			return
		}

		// Validasi radius hanya untuk absen mandiri siswa.
		// Take over oleh guru/DUDI = entri manual terpercaya (lokasi bisa di luar DUDI).
		if !takeOver && student.DudiID != nil {
			var dudi models.DUDI
			if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
				distance := haversine(req.Latitude, req.Longitude, dudi.Latitude, dudi.Longitude)
				if distance > float64(dudi.RadiusAllowed) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":      "Anda berada di luar radius DUDI",
						"distance_m": distance,
						"radius_m":   dudi.RadiusAllowed,
					})
					return
				}
			}
		}

		ip, ua := clientDevice(c, takeOver)

		absensi := models.Absensi{
			StudentID:  sid,
			Timestamp:  now,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			Type:       "masuk",
			Status:     "hadir",
			IsVerified: true,
			IPAddress:  ip,
			UserAgent:  ua,
		}

		if err := database.DB.Create(&absensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
			return
		}

		// Deteksi perangkat bersama hanya untuk absen mandiri siswa
		// (take over memang 1 perangkat untuk banyak siswa).
		if !takeOver {
			flagSharedDevice(&absensi, todayStart, todayEnd)
		}

		c.JSON(http.StatusCreated, gin.H{
			"message": "Absen masuk tercatat",
			"data": gin.H{
				"id":          absensi.ID,
				"type":        absensi.Type,
				"timestamp":   absensi.Timestamp,
				"status":      absensi.Status,
				"is_verified": absensi.IsVerified,
			},
		})
		return
	}

	// req.Type == "pulang"

	// Cek sudah absen masuk hari ini?
	var masuk models.Absensi
	if err := database.DB.Where("student_id = ? AND type = 'masuk' AND timestamp >= ? AND timestamp < ?", sid, todayStart, todayEnd).First(&masuk).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": subject + " belum absen masuk hari ini"})
		return
	}

	// Cek sudah absen pulang hari ini?
	var count int64
	database.DB.Model(&models.Absensi{}).
		Where("student_id = ? AND type = 'pulang' AND timestamp >= ? AND timestamp < ?", sid, todayStart, todayEnd).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": subject + " sudah absen pulang hari ini"})
		return
	}

	// Minimal 6 jam setelah absen masuk.
	// Take over oleh guru/DUDI dikecualikan (mengoreksi absen yang terlewat).
	if !takeOver {
		checkoutAvailable := masuk.Timestamp.Add(6 * time.Hour)
		if now.Before(checkoutAvailable) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":            "Absen pulang belum tersedia",
				"available_at":     checkoutAvailable,
				"available_at_wib": checkoutAvailable.Format("15:04 WIB"),
			})
			return
		}
	}

	// Validasi radius hanya untuk absen mandiri siswa.
	if !takeOver && student.DudiID != nil {
		var dudi models.DUDI
		if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
			distance := haversine(req.Latitude, req.Longitude, dudi.Latitude, dudi.Longitude)
			if distance > float64(dudi.RadiusAllowed) {
				c.JSON(http.StatusBadRequest, gin.H{
					"error":      "Anda berada di luar radius DUDI",
					"distance_m": distance,
					"radius_m":   dudi.RadiusAllowed,
				})
				return
			}
		}
	}

	ip, ua := clientDevice(c, takeOver)

	absensi := models.Absensi{
		StudentID:  sid,
		Timestamp:  now,
		Latitude:   req.Latitude,
		Longitude:  req.Longitude,
		Type:       "pulang",
		Status:     "hadir",
		IsVerified: takeOver,
		IPAddress:  ip,
		UserAgent:  ua,
	}

	if err := database.DB.Create(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
		return
	}

	if !takeOver {
		flagSharedDevice(&absensi, todayStart, todayEnd)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Absen pulang tercatat",
		"data": gin.H{
			"id":          absensi.ID,
			"type":        absensi.Type,
			"timestamp":   absensi.Timestamp,
			"status":      absensi.Status,
			"is_verified": absensi.IsVerified,
		},
	})
}

func (h *AbsensiHandler) History(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")
	jurusan, _ := c.Get("jurusan")

	uid, _ := uuid.Parse(userID.(string))

	var absensiList []models.Absensi

	query := database.DB.Order("timestamp DESC")

	if role == "student" {
		query = query.Where("student_id = ?", uid)
	} else {
		studentID := c.Query("student_id")
		jurusanFilter := c.Query("jurusan")

		if role == "teacher" {
			query = query.Joins("JOIN users ON users.id = absensis.student_id").Where("users.teacher_id = ?", uid)
		} else if role == "admin_jurusan" && jurusan != nil && jurusan.(string) != "" {
			query = query.Joins("JOIN users ON users.id = absensis.student_id").Where("users.jurusan = ?", jurusan.(string))
		} else if jurusanFilter != "" {
			query = query.Joins("JOIN users ON users.id = absensis.student_id").Where("users.jurusan = ?", jurusanFilter)
		}

		if studentID != "" {
			query = query.Where("absensis.student_id = ?", studentID)
		}
	}

	query.Find(&absensiList)

	if role != "student" && len(absensiList) > 0 {
		var studentIDs []uuid.UUID
		for _, a := range absensiList {
			studentIDs = append(studentIDs, a.StudentID)
		}
		var students []models.User
		database.DB.Where("id IN ?", studentIDs).Find(&students)
		studentMap := make(map[uuid.UUID]*models.User)
		for i := range students {
			studentMap[students[i].ID] = &students[i]
		}
		for i := range absensiList {
			absensiList[i].Student = studentMap[absensiList[i].StudentID]
		}
	}

	c.JSON(http.StatusOK, gin.H{"data": absensiList})
}

func (h *AbsensiHandler) Status(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students"})
		return
	}

	uid, _ := uuid.Parse(userID.(string))

	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	var masuk models.Absensi
	hasMasuk := database.DB.Where("student_id = ? AND type = 'masuk' AND timestamp >= ? AND timestamp < ?", uid, todayStart, todayEnd).First(&masuk).Error == nil

	var pulang models.Absensi
	hasPulang := database.DB.Where("student_id = ? AND type = 'pulang' AND timestamp >= ? AND timestamp < ?", uid, todayStart, todayEnd).First(&pulang).Error == nil

	result := gin.H{
		"has_masuk":  hasMasuk,
		"has_pulang": hasPulang,
		"server_time": now.Format("2006-01-02 15:04:05 WIB"),
	}

	if hasMasuk {
		result["masuk_at"] = masuk.Timestamp
		checkoutAvailable := masuk.Timestamp.Add(6 * time.Hour)
		result["pulang_available_at"] = checkoutAvailable
		result["pulang_available"] = now.After(checkoutAvailable) || now.Equal(checkoutAvailable)
		// Display times in WIB
		result["masuk_at_wib"] = masuk.Timestamp.In(loc).Format("15:04 WIB")
		result["pulang_available_at_wib"] = checkoutAvailable.In(loc).Format("15:04 WIB")
	}

	if hasPulang {
		result["pulang_at"] = pulang.Timestamp
	}

	c.JSON(http.StatusOK, result)
}

func (h *AbsensiHandler) Verify(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers/admin can verify attendance"})
		return
	}

	id := c.Param("id")
	absensiID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var absensi models.Absensi
	if err := database.DB.First(&absensi, "id = ?", absensiID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attendance not found"})
		return
	}

	absensi.IsVerified = true
	database.DB.Save(&absensi)

	c.JSON(http.StatusOK, gin.H{"message": "Attendance verified", "data": absensi})
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371000
	dLat := (lat2 - lat1) * (math.Pi / 180.0)
	dLon := (lon2 - lon1) * (math.Pi / 180.0)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180.0))*math.Cos(lat2*(math.Pi/180.0))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func flagSharedDevice(record *models.Absensi, start, end time.Time) {
	if record.IPAddress == "" {
		return
	}

	var others []models.Absensi
	database.DB.
		Where("student_id != ? AND ip_address = ? AND user_agent = ? AND timestamp >= ? AND timestamp < ?",
			record.StudentID, record.IPAddress, record.UserAgent, start, end).
		Find(&others)

	if len(others) > 0 {
		record.IsSuspicious = true
		database.DB.Model(record).Update("is_suspicious", true)

		for _, o := range others {
			if !o.IsSuspicious {
				o.IsSuspicious = true
				database.DB.Model(&o).Update("is_suspicious", true)
			}
		}
	}
}

type IzinRequest struct {
	StudentID string `json:"student_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
	Tanggal   string `json:"tanggal"`
}

func (h *AbsensiHandler) CreateIzin(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "teacher" && role != "dudi" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya guru dan instruktur DUDI yang dapat mencatat izin/sakit"})
		return
	}

	var req IzinRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "student_id dan status diperlukan"})
		return
	}

	if req.Status != "izin" && req.Status != "sakit" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status harus 'izin' atau 'sakit'"})
		return
	}

	uid, _ := uuid.Parse(userID.(string))

	student, msg := resolveStudentForStaff(uid, role.(string), req.StudentID)
	if student == nil {
		status := http.StatusForbidden
		if msg == "Siswa tidak ditemukan" || msg == "ID siswa tidak valid" {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": msg})
		return
	}

	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	targetDate := now
	if req.Tanggal != "" {
		parsed, err := time.ParseInLocation("2006-01-02", req.Tanggal, loc)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal tidak valid, gunakan YYYY-MM-DD"})
			return
		}
		parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 8, 0, 0, 0, loc)
		targetDate = parsed
	} else {
		targetDate = time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, loc)
	}

	dayStart := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, loc)
	dayEnd := dayStart.Add(24 * time.Hour)

	var existingCount int64
	database.DB.Model(&models.Absensi{}).
		Where("student_id = ? AND timestamp >= ? AND timestamp < ?", student.ID, dayStart, dayEnd).
		Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Siswa sudah memiliki catatan absensi pada tanggal ini"})
		return
	}

	ip, ua := clientDevice(c, true)

	absensi := models.Absensi{
		StudentID:  student.ID,
		Timestamp:  targetDate,
		Latitude:   0,
		Longitude:  0,
		Type:       "masuk",
		Status:     req.Status,
		IsVerified: true,
		IPAddress:  ip,
		UserAgent:  ua,
	}

	if err := database.DB.Create(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Absensi " + req.Status + " tercatat",
		"data": gin.H{
			"id":          absensi.ID,
			"student_id":  absensi.StudentID,
			"type":        absensi.Type,
			"status":      absensi.Status,
			"timestamp":   absensi.Timestamp,
			"is_verified": absensi.IsVerified,
		},
	})
}

func (h *AbsensiHandler) SuspiciousReport(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	uid, _ := uuid.Parse(userID.(string))

	loc := time.FixedZone("WIB", 7*3600)
	now := time.Now().In(loc)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	todayEnd := todayStart.Add(24 * time.Hour)

	type SuspiciousEntry struct {
		ID        uuid.UUID `json:"id"`
		StudentID uuid.UUID `json:"student_id"`
		Name      string    `json:"student_name"`
		Type      string    `json:"type"`
		Timestamp time.Time `json:"timestamp"`
		IPAddress string    `json:"ip_address"`
		UserAgent string    `json:"user_agent"`
		ISVerified bool     `json:"is_verified"`
	}

	var entries []SuspiciousEntry
	query := database.DB.Table("absensis").
		Select("absensis.id, absensis.student_id, users.full_name as name, absensis.type, absensis.timestamp, absensis.ip_address, absensis.user_agent, absensis.is_verified").
		Joins("JOIN users ON users.id = absensis.student_id").
		Where("absensis.is_suspicious = true AND absensis.timestamp >= ? AND absensis.timestamp < ?", todayStart, todayEnd).
		Order("absensis.timestamp DESC")

	if role == "teacher" {
		query = query.Where("users.teacher_id = ?", uid)
	} else if role == "dudi" {
		var dudiUser models.User
		if database.DB.First(&dudiUser, "id = ?", uid).Error == nil && dudiUser.DudiID != nil {
			query = query.Where("users.dudi_id = ?", dudiUser.DudiID)
		}
	}

	query.Scan(&entries)

	c.JSON(http.StatusOK, gin.H{"data": entries})
}
