package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type ReportHandler struct{}

func NewReportHandler() *ReportHandler {
	return &ReportHandler{}
}

func withPeriodeFilter(query *gorm.DB, periodeID string) (*gorm.DB, bool) {
	if periodeID == "" {
		var active models.Periode
		if database.DB.Where("is_active = ?", true).First(&active).Error == nil {
			return query.Where("timestamp >= ? AND timestamp < ?::date + interval '1 day'",
				active.StartDate, active.EndDate), true
		}
		return query, false
	}

	id, err := uuid.Parse(periodeID)
	if err != nil {
		return query, false
	}

	var periode models.Periode
	if database.DB.First(&periode, "id = ?", id).Error != nil {
		return query, false
	}

	return query.Where("timestamp >= ? AND timestamp < ?::date + interval '1 day'",
		periode.StartDate, periode.EndDate), true
}

func (h *ReportHandler) AbsensiReport(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers and admin can access reports"})
		return
	}

	userID, _ := c.Get("user_id")
	periodeID := c.Query("periode_id")
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

	query, _ = withPeriodeFilter(query, periodeID)

	query.Find(&absensiList)

	summary := make(map[string]interface{})
	var totalHadir, totalTerlambat, totalIzin, totalSakit int64

	summaryQuery := database.DB.Model(&models.Absensi{})
	if studentID != "" {
		summaryQuery = summaryQuery.Where("student_id = ?", studentID)
	}
	summaryQuery, _ = withPeriodeFilter(summaryQuery, periodeID)

	summaryQuery.Where("status = 'hadir'").Count(&totalHadir)
	summaryQuery.Where("status = 'terlambat'").Count(&totalTerlambat)
	summaryQuery.Where("status = 'izin'").Count(&totalIzin)
	summaryQuery.Where("status = 'sakit'").Count(&totalSakit)

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
