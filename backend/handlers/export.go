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
	jurusan, _ := c.Get("jurusan")
	if role != "teacher" && role != "admin" && role != "admin_jurusan" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	periodeID := c.Query("periode_id")
	studentID := c.Query("student_id")
	jurusanFilter := c.Query("jurusan")

	var absensiList []models.Absensi
	query := database.DB.Preload("Student").Joins("JOIN users ON users.id = absensis.student_id").Order("timestamp DESC")

	if role == "admin_jurusan" && jurusan != nil && jurusan.(string) != "" {
		query = query.Where("users.jurusan = ?", jurusan.(string))
	} else if jurusanFilter != "" {
		query = query.Where("users.jurusan = ?", jurusanFilter)
	}

	if studentID != "" {
		query = query.Where("absensis.student_id = ?", studentID)
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
	jurusan, _ := c.Get("jurusan")
	if role != "teacher" && role != "admin" && role != "admin_jurusan" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	studentID := c.Query("student_id")
	jurusanFilter := c.Query("jurusan")

	var jurnalList []models.Jurnal
	query := database.DB.Preload("Student").Joins("JOIN users ON users.id = jurnals.student_id").Order("date DESC")

	if role == "admin_jurusan" && jurusan != nil && jurusan.(string) != "" {
		query = query.Where("users.jurusan = ?", jurusan.(string))
	} else if jurusanFilter != "" {
		query = query.Where("users.jurusan = ?", jurusanFilter)
	}

	if studentID != "" {
		query = query.Where("jurnals.student_id = ?", studentID)
	}

	query.Find(&jurnalList)

	filename := fmt.Sprintf("jurnal_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Tanggal", "Nama Siswa", "NIS", "Kegiatan", "Refleksi", "Komentar Guru", "Komentar Instruktur"})

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
	jurusan, _ := c.Get("jurusan")
	if role != "teacher" && role != "admin" && role != "admin_jurusan" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	jurusanFilter := c.Query("jurusan")

	query := database.DB.Preload("Student").Preload("Student.DUDI")
	if role == "admin_jurusan" && jurusan != nil && jurusan.(string) != "" {
		query = query.Joins("JOIN users ON users.id = penilaians.student_id").Where("users.jurusan = ?", jurusan.(string))
	} else if jurusanFilter != "" {
		query = query.Joins("JOIN users ON users.id = penilaians.student_id").Where("users.jurusan = ?", jurusanFilter)
	}

	var penilaianList []models.Penilaian
	query.Order("final_score DESC").Find(&penilaianList)

	filename := fmt.Sprintf("nilai_%s.csv", time.Now().Format("20060102_150405"))
	c.Header("Content-Type", "text/csv; charset=utf-8")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	writer := csv.NewWriter(c.Writer)
	writer.Write([]string{"Nama Siswa", "NIS", "DUDI", "Kehadiran (%)", "Alur Bisnis (1-5)", "Soft Skills (1-5)", "Kompetensi Teknis (1-5)", "POS & K3LH (1-5)", "Nilai Akhir", "Grade", "Catatan"})

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
