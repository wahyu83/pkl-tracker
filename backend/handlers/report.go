package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

func (h *ReportHandler) AbsensiReport(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers and admin can access reports"})
		return
	}

	userID, _ := c.Get("user_id")
	period := c.Query("periode")
	studentID := c.Query("student_id")

	var absensiList []models.Absensi
	query := database.DB.Preload("Student").Order("timestamp DESC")

	if role == "teacher" {
		uid, _ := uuid.Parse(userID.(string))
		query = query.Joins("JOIN users ON users.id = absensis.student_id").
			Where("users.teacher_id = ?", uid)
	}

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}

	if period != "" {
		query = query.Where("DATE(timestamp) BETWEEN ? AND ?", period+"-01", period+"-31")
	}

	query.Find(&absensiList)

	summary := make(map[string]interface{})
	var totalHadir, totalTerlambat, totalIzin, totalSakit int64

	database.DB.Model(&models.Absensi{}).Where("status = 'hadir'").Count(&totalHadir)
	database.DB.Model(&models.Absensi{}).Where("status = 'terlambat'").Count(&totalTerlambat)
	database.DB.Model(&models.Absensi{}).Where("status = 'izin'").Count(&totalIzin)
	database.DB.Model(&models.Absensi{}).Where("status = 'sakit'").Count(&totalSakit)

	summary["total_hadir"] = totalHadir
	summary["total_terlambat"] = totalTerlambat
	summary["total_izin"] = totalIzin
	summary["total_sakit"] = totalSakit
	summary["total_records"] = len(absensiList)

	c.JSON(http.StatusOK, gin.H{
		"summary": summary,
		"data":    absensiList,
	})
}

func (h *ReportHandler) JurnalReport(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers and admin can access reports"})
		return
	}

	studentID := c.Query("student_id")

	var jurnalList []models.Jurnal
	query := database.DB.Preload("Student").Order("date DESC")

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}

	query.Find(&jurnalList)

	c.JSON(http.StatusOK, gin.H{
		"total": len(jurnalList),
		"data":  jurnalList,
	})
}

func (h *ReportHandler) NilaiReport(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers and admin can access reports"})
		return
	}

	var penilaianList []models.Penilaian
	database.DB.Preload("Student").Order("final_score DESC").Find(&penilaianList)

	c.JSON(http.StatusOK, gin.H{
		"total": len(penilaianList),
		"data":  penilaianList,
	})
}
