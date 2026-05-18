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

		isVerified := false
		if student.DudiID != nil {
			var dudi models.DUDI
			if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
				distance := haversine(req.Latitude, req.Longitude, dudi.Latitude, dudi.Longitude)
				isVerified = distance <= float64(dudi.RadiusAllowed)
			}
		}

		absensi := models.Absensi{
			StudentID:  uid,
			Timestamp:  now,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			Type:       "masuk",
			Status:     status,
			IsVerified: isVerified,
		}

		if err := database.DB.Create(&absensi).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan absensi"})
			return
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

		// Cek 7 jam setelah absen masuk
		checkoutAvailable := masuk.Timestamp.Add(7 * time.Hour)
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
		}

		if student.DudiID != nil {
			var dudi models.DUDI
			if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
				distance := haversine(req.Latitude, req.Longitude, dudi.Latitude, dudi.Longitude)
				if distance > float64(dudi.RadiusAllowed) {
					c.JSON(http.StatusBadRequest, gin.H{
						"error":        "Anda berada di luar area DUDI",
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
		checkoutAvailable := masuk.Timestamp.Add(7 * time.Hour)
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
