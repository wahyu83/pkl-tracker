package handlers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type ExportHandler struct{}

func NewExportHandler() *ExportHandler {
	return &ExportHandler{}
}

func (h *ExportHandler) ExportAbsensi(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	periodeID := c.Query("periode_id")
	studentID := c.Query("student_id")

	var absensiList []models.Absensi
	query := database.DB.Preload("Student").Order("timestamp DESC")

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}
	query, _ = withPeriodeFilter(query, periodeID)

	query.Find(&absensiList)

	filename := fmt.Sprintf("absensi_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Tanggal", "Nama Siswa", "NIS", "DUDI", "Status", "Latitude", "Longitude", "Terverifikasi", "Foto"})

	for _, a := range absensiList {
		studentName := ""
		nis := ""
		dudi := ""
		if a.Student != nil {
			studentName = a.Student.FullName
			nis = a.Student.NisNipNik
			if a.Student.DUDI != nil {
				dudi = a.Student.DUDI.CompanyName
			}
		}
		verified := "Tidak"
		if a.IsVerified {
			verified = "Ya"
		}
		status := statusLabel(a.Status)
		writer.Write([]string{
			a.Timestamp.Format("2006-01-02 15:04:05"),
			studentName,
			nis,
			dudi,
			status,
			fmt.Sprintf("%.6f", a.Latitude),
			fmt.Sprintf("%.6f", a.Longitude),
			verified,
			a.PhotoURL,
		})
	}
	writer.Flush()
}

func (h *ExportHandler) ExportJurnal(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	studentID := c.Query("student_id")

	var jurnalList []models.Jurnal
	query := database.DB.Preload("Student").Order("date DESC")

	if studentID != "" {
		query = query.Where("student_id = ?", studentID)
	}

	query.Find(&jurnalList)

	filename := fmt.Sprintf("jurnal_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Tanggal", "Nama Siswa", "NIS", "Kegiatan", "Refleksi", "Komentar Guru", "Komentar DUDI"})

	for _, j := range jurnalList {
		studentName := ""
		nis := ""
		if j.Student != nil {
			studentName = j.Student.FullName
			nis = j.Student.NisNipNik
		}
		writer.Write([]string{
			j.Date.Format("2006-01-02"),
			studentName,
			nis,
			j.Activity,
			j.Reflection,
			j.TeacherComment,
			j.DudiComment,
		})
	}
	writer.Flush()
}

func (h *ExportHandler) ExportNilai(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "teacher" && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var penilaianList []models.Penilaian
	database.DB.Preload("Student").Preload("Student.DUDI").Order("final_score DESC").Find(&penilaianList)

	filename := fmt.Sprintf("nilai_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Nama Siswa", "NIS", "DUDI", "Kehadiran (%)", "Disiplin (1-5)", "Tanggung Jawab (1-5)", "Kerjasama (1-5)", "Inisiatif (1-5)", "Nilai Akhir", "Grade", "Catatan"})

	for _, p := range penilaianList {
		studentName := ""
		nis := ""
		dudi := ""
		if p.Student != nil {
			studentName = p.Student.FullName
			nis = p.Student.NisNipNik
			if p.Student.DUDI != nil {
				dudi = p.Student.DUDI.CompanyName
			}
		}
		writer.Write([]string{
			studentName,
			nis,
			dudi,
			fmt.Sprintf("%.2f", p.AttendanceScoreAuto),
			fmt.Sprintf("%d", p.Discipline),
			fmt.Sprintf("%d", p.Responsibility),
			fmt.Sprintf("%d", p.Teamwork),
			fmt.Sprintf("%d", p.Initiative),
			fmt.Sprintf("%.2f", p.FinalScore),
			p.FinalGrade,
			p.Notes,
		})
	}
	writer.Flush()
}

func statusLabel(s string) string {
	m := map[string]string{
		"hadir":     "Hadir",
		"terlambat": "Terlambat",
		"izin":      "Izin",
		"sakit":     "Sakit",
	}
	if v, ok := m[s]; ok {
		return v
	}
	return s
}
