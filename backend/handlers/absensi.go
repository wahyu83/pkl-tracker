package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
	"pkl-tracker/storage"
)

type AbsensiHandler struct {
	store *storage.DriveStorage
}

func NewAbsensiHandler(store *storage.DriveStorage) *AbsensiHandler {
	return &AbsensiHandler{store: store}
}

type AbsensiRequest struct {
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

	latStr := c.PostForm("latitude")
	lngStr := c.PostForm("longitude")
	status := c.PostForm("status")

	var lat, lng float64
	fmt.Sscanf(latStr, "%f", &lat)
	fmt.Sscanf(lngStr, "%f", &lng)

	if lat == 0 && lng == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Latitude and longitude required"})
		return
	}

	if status == "" {
		status = "hadir"
		now := time.Now().In(time.FixedZone("WIB", 7*3600))
		cutoff := time.Date(now.Year(), now.Month(), now.Day(), 8, 0, 0, 0, now.Location())
		if now.After(cutoff) {
			status = "terlambat"
		}
	}

	uid, _ := uuid.Parse(userID.(string))

	var student models.User
	if err := database.DB.First(&student, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	isVerified := false
	if student.DudiID != nil {
		var dudi models.DUDI
		if err := database.DB.First(&dudi, "id = ?", *student.DudiID).Error; err == nil {
			distance := haversine(lat, lng, dudi.Latitude, dudi.Longitude)
			isVerified = distance <= float64(dudi.RadiusAllowed)
		}
	}

	photoURL := ""

	file, header, err := c.Request.FormFile("photo")
	if err == nil {
		defer file.Close()
		timestamp := time.Now().Unix()
		ext := ".jpg"
		if strings.Contains(header.Filename, ".") {
			parts := strings.Split(header.Filename, ".")
			ext = "." + parts[len(parts)-1]
		}
		filename := fmt.Sprintf("absensi_%d%s", timestamp, ext)

		uploadedURL, uploadErr := h.store.UploadFile(file, filename, uid.String())
		if uploadErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload photo: " + uploadErr.Error()})
			return
		}
		photoURL = uploadedURL
	}

	absensi := models.Absensi{
		StudentID:  uid,
		Timestamp:  time.Now(),
		Latitude:   lat,
		Longitude:  lng,
		PhotoURL:   photoURL,
		Status:     status,
		IsVerified: isVerified,
	}

	if err := database.DB.Create(&absensi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attendance"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Attendance recorded",
		"data": gin.H{
			"id":          absensi.ID,
			"timestamp":   absensi.Timestamp,
			"status":      absensi.Status,
			"is_verified": absensi.IsVerified,
			"photo_url":   absensi.PhotoURL,
		},
	})
}

func (h *AbsensiHandler) History(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	uid, _ := uuid.Parse(userID.(string))

	var absensiList []models.Absensi

	query := database.DB.Order("timestamp DESC")

	if role == "student" {
		query = query.Where("student_id = ?", uid)
	} else {
		studentID := c.Query("student_id")
		if studentID != "" {
			query = query.Where("student_id = ?", studentID)
		}
	}

	query.Find(&absensiList)

	c.JSON(http.StatusOK, gin.H{"data": absensiList})
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
