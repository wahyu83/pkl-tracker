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
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	Status    string  `json:"status"`
}

func (h *AbsensiHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can submit attendance"})
		return
	}

	var req AbsensiRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude required"})
		return
	}

	if req.Latitude == 0 && req.Longitude == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude required"})
		return
	}

	if req.Type != "masuk" && req.Type != "pulang" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Type must be 'masuk' or 'pulang'"})
		return
	}

	uid, _ := uuid.Parse(userID.(string))

	var student models.User
	if err := database.DB.First(&student, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Gunakan waktu server (WIB)
	now := time.Now().In(time.FixedZone("WIB", 7*3600))
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)

	if req.Type == "masuk" {
		// Cek sudah absen masuk hari ini?
		var count int64
		database.DB.Model(&models.Absensi{}).
			Where("student_id = ? AND type = 'masuk' AND timestamp >= ? AND timestamp < ?", uid, todayStart, todayEnd).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah absen masuk hari ini"})
			return
		}

		status := "hadir"

		if student.DudiID != nil {
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

		isVerified := true

		absensi := models.Absensi{
			StudentID:  uid,
			Timestamp:  now,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			Type:       "masuk",
			Status:     status,
			IsVerified: isVerified,
			IPAddress:  c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		if err := database.DB.Create(&absensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
			return
		}

		flagSharedDevice(&absensi, todayStart, todayEnd)

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

	if req.Type == "pulang" {
		// Cek sudah absen masuk hari ini?
		var masuk models.Absensi
		if err := database.DB.Where("student_id = ? AND type = 'masuk' AND timestamp >= ? AND timestamp < ?", uid, todayStart, todayEnd).First(&masuk).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum absen masuk hari ini"})
			return
		}

		// Cek sudah absen pulang hari ini?
		var count int64
		database.DB.Model(&models.Absensi{}).
			Where("student_id = ? AND type = 'pulang' AND timestamp >= ? AND timestamp < ?", uid, todayStart, todayEnd).
			Count(&count)
		if count > 0 {
			c.JSON(http.StatusConflict, gin.H{"error": "Anda sudah absen pulang hari ini"})
			return
		}

		// Cek 6 jam setelah absen masuk
		checkoutAvailable := masuk.Timestamp.Add(6 * time.Hour)
		if now.Before(checkoutAvailable) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":             "Absen pulang belum tersedia",
				"available_at":      checkoutAvailable,
				"available_at_wib":  checkoutAvailable.Format("15:04 WIB"),
			})
			return
		}

		absensi := models.Absensi{
			StudentID: uid,
			Timestamp: now,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
			Type:      "pulang",
			Status:    "hadir",
			IPAddress: c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
		}

		if student.DudiID != nil {
			var dudi models.DUDI
			if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
				distance := haversine(req.Latitude, req.Longitude, dudi.Latitude, dudi.Longitude)
				if distance > float64(dudi.RadiusAllowed) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":      "Anda berada di luar radius DUDI",
						"distance_m":   distance,
						"radius_m":     dudi.RadiusAllowed,
					})
					return
				}
			}
		}

		if err := database.DB.Create(&absensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
			return
		}

		flagSharedDevice(&absensi, todayStart, todayEnd)

		c.JSON(http.StatusCreated, gin.H{
			"message": "Absen pulang tercatat",
			"data": gin.H{
				"id":        absensi.ID,
				"type":      absensi.Type,
				"timestamp": absensi.Timestamp,
				"status":    absensi.Status,
			},
		})
		return
	}
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
	studentUID, err := uuid.Parse(req.StudentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID siswa tidak valid"})
		return
	}

	var student models.User
	if err := database.DB.First(&student, "id = ? AND role = 'student'", studentUID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan"})
		return
	}

	if role == "teacher" {
		if student.TeacherID == nil || *student.TeacherID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Siswa ini tidak dalam bimbingan Anda"})
			return
		}
	} else if role == "dudi" {
		var dudiUser models.User
		if database.DB.First(&dudiUser, "id = ?", uid).Error != nil || dudiUser.DudiID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Profil DUDI tidak ditemukan"})
			return
		}
		if student.DudiID == nil || *student.DudiID != *dudiUser.DudiID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Siswa ini tidak berada di DUDI Anda"})
			return
		}
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
		Where("student_id = ? AND timestamp >= ? AND timestamp < ?", studentUID, dayStart, dayEnd).
		Count(&existingCount)
	if existingCount > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Siswa sudah memiliki catatan absensi pada tanggal ini"})
		return
	}

	absensi := models.Absensi{
		StudentID:  studentUID,
		Timestamp:  targetDate,
		Latitude:   0,
		Longitude:  0,
		Type:       "masuk",
		Status:     req.Status,
		IsVerified: true,
		IPAddress:  c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
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
