package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type PenilaianHandler struct{}

func NewPenilaianHandler() *PenilaianHandler {
	return &PenilaianHandler{}
}

type PenilaianRequest struct {
	StudentID           string `json:"student_id" binding:"required"`
	Discipline          int    `json:"discipline" binding:"min=0,max=100"`
	Responsibility      int    `json:"responsibility" binding:"min=0,max=100"`
	Teamwork            int    `json:"teamwork" binding:"min=0,max=100"`
	Initiative          int    `json:"initiative" binding:"min=0,max=100"`
	Notes               string `json:"notes"`
	AttendanceStartDate string `json:"attendance_start_date"`
	AttendanceEndDate   string `json:"attendance_end_date"`
}

// parseDateParam membaca tanggal "YYYY-MM-DD" (WIB). String kosong -> nil.
func parseDateParam(value string) (*time.Time, error) {
	if strings.TrimSpace(value) == "" {
		return nil, nil
	}
	loc := time.FixedZone("WIB", 7*3600)
	t, err := time.ParseInLocation("2006-01-02", value, loc)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (h *PenilaianHandler) CreateOrUpdate(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "dudi" && role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya guru dan instruktur DUDI yang dapat mengisi nilai"})
		return
	}

	var req PenilaianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	uid, _ := uuid.Parse(userID.(string))

	var student models.User
	if err := database.DB.First(&student, "id = ? AND role = 'student'", studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	// Hak akses: guru hanya siswa bimbingannya, DUDI hanya siswa di DUDI-nya.
	if role == "teacher" {
		if student.TeacherID == nil || *student.TeacherID != uid {
			c.JSON(http.StatusForbidden, gin.H{"error": "Siswa ini tidak dalam bimbingan Anda"})
			return
		}
	} else {
		var dudiUser models.User
		if err := database.DB.First(&dudiUser, "id = ?", uid).Error; err != nil || dudiUser.DudiID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your account is not linked to a company"})
			return
		}
		if student.DudiID == nil || *student.DudiID != *dudiUser.DudiID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Student is not assigned to your company"})
			return
		}
	}

	if student.DudiID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Siswa belum memiliki DUDI"})
		return
	}

	startDate, err := parseDateParam(req.AttendanceStartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal mulai tidak valid (YYYY-MM-DD)"})
		return
	}
	endDate, err := parseDateParam(req.AttendanceEndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format tanggal selesai tidak valid (YYYY-MM-DD)"})
		return
	}
	if startDate != nil && endDate != nil && endDate.Before(*startDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tanggal selesai tidak boleh sebelum tanggal mulai"})
		return
	}

	var penilaian models.Penilaian
	result := database.DB.Where("student_id = ? AND dudi_id = ?", studentID, *student.DudiID).First(&penilaian)
	isNew := result.Error != nil
	if isNew {
		penilaian = models.Penilaian{
			StudentID: studentID,
			DudiID:    *student.DudiID,
		}
	}

	penilaian.Discipline = req.Discipline
	penilaian.Responsibility = req.Responsibility
	penilaian.Teamwork = req.Teamwork
	penilaian.Initiative = req.Initiative
	penilaian.Notes = req.Notes
	penilaian.SubmittedAt = time.Now()

	// Rentang tanggal pelaksanaan PKL hanya ditentukan oleh guru.
	if role == "teacher" {
		penilaian.AttendanceStartDate = startDate
		penilaian.AttendanceEndDate = endDate
	}

	penilaian.AttendanceScoreAuto = calculateAttendanceScore(
		studentID, penilaian.AttendanceStartDate, penilaian.AttendanceEndDate)
	penilaian.FinalScore, penilaian.FinalGrade = calculateFinalScore(penilaian)

	if isNew {
		database.DB.Create(&penilaian)
	} else {
		database.DB.Save(&penilaian)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Evaluation saved", "data": penilaian})
}

func (h *PenilaianHandler) GetByStudent(c *gin.Context) {
	studentID := c.Param("studentId")
	sID, err := uuid.Parse(studentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	var penilaian models.Penilaian
	if err := database.DB.Where("student_id = ?", sID).First(&penilaian).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evaluation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": penilaian})
}

func (h *PenilaianHandler) List(c *gin.Context) {
	role, _ := c.Get("role")
	jurusan, _ := c.Get("jurusan")

	var penilaianList []models.Penilaian
	query := database.DB.Preload("Student").Preload("Student.DUDI")

	switch role {
	case "dudi":
		userID, _ := c.Get("user_id")
		uid, _ := uuid.Parse(userID.(string))
		query = query.Where("dudi_id = ?", uid)
	case "teacher":
		userID, _ := c.Get("user_id")
		uid, _ := uuid.Parse(userID.(string))
		query = query.Joins("JOIN users ON users.id = penilaians.student_id").Where("users.teacher_id = ?", uid)
	case "admin", "admin_jurusan":
		studentID := c.Query("student_id")
		jurusanFilter := c.Query("jurusan")

		if role == "admin_jurusan" && jurusan != nil && jurusan.(string) != "" {
			query = query.Joins("JOIN users ON users.id = penilaians.student_id").Where("users.jurusan = ?", jurusan.(string))
		} else if jurusanFilter != "" {
			query = query.Joins("JOIN users ON users.id = penilaians.student_id").Where("users.jurusan = ?", jurusanFilter)
		}

		if studentID != "" {
			query = query.Where("penilaians.student_id = ?", studentID)
		}
	default:
		userID, _ := c.Get("user_id")
		uid, _ := uuid.Parse(userID.(string))
		query = query.Where("student_id = ?", uid)
	}

	query.Find(&penilaianList)
	c.JSON(http.StatusOK, gin.H{"data": penilaianList})
}

// attendanceScope membatasi absensi ke satu siswa, dan ke rentang tanggal
// pelaksanaan PKL bila rentangnya diisi.
func attendanceScope(studentID uuid.UUID, start, end *time.Time) *gorm.DB {
	query := database.DB.Model(&models.Absensi{}).Where("student_id = ?", studentID)
	if start != nil && end != nil {
		query = query.Where("timestamp >= ? AND timestamp < ?", *start, end.AddDate(0, 0, 1))
	}
	return query
}

func calculateAttendanceScore(studentID uuid.UUID, start, end *time.Time) float64 {
	var total int64
	attendanceScope(studentID, start, end).Count(&total)
	if total == 0 {
		return 0
	}

	var present int64
	attendanceScope(studentID, start, end).
		Where("status IN ('hadir','terlambat')").
		Count(&present)

	return float64(present) / float64(total) * 100
}

func calculateFinalScore(p models.Penilaian) (float64, string) {
	attendanceWeight := 0.3
	manualWeight := 0.7

	manualAvg := float64(p.Discipline+p.Responsibility+p.Teamwork+p.Initiative) / 4.0
	score := (p.AttendanceScoreAuto * attendanceWeight) + (manualAvg * manualWeight)

	var grade string
	switch {
	case score >= 90:
		grade = "A"
	case score >= 80:
		grade = "B+"
	case score >= 70:
		grade = "B"
	case score >= 60:
		grade = "C"
	default:
		grade = "D"
	}

	return score, grade
}
