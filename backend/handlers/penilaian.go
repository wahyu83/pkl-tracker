package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type PenilaianHandler struct{}

func NewPenilaianHandler() *PenilaianHandler {
	return &PenilaianHandler{}
}

type PenilaianRequest struct {
	StudentID      string  `json:"student_id" binding:"required"`
	Discipline     int     `json:"discipline" binding:"required,min=1,max=5"`
	Responsibility int     `json:"responsibility" binding:"required,min=1,max=5"`
	Teamwork       int     `json:"teamwork" binding:"required,min=1,max=5"`
	Initiative     int     `json:"initiative" binding:"required,min=1,max=5"`
	Notes          string  `json:"notes"`
}

func (h *PenilaianHandler) CreateOrUpdate(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "dudi" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only DUDI can submit evaluations"})
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
	if err := database.DB.First(&student, "id = ?", studentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}

	if student.DudiID == nil || *student.DudiID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Student is not assigned to your company"})
		return
	}

	attendanceScore := calculateAttendanceScore(studentID)

	var penilaian models.Penilaian
	result := database.DB.Where("student_id = ? AND dudi_id = ?", studentID, uid).First(&penilaian)

	if result.Error != nil {
		penilaian = models.Penilaian{
			StudentID:           studentID,
			DudiID:              uid,
			AttendanceScoreAuto: attendanceScore,
			Discipline:          req.Discipline,
			Responsibility:      req.Responsibility,
			Teamwork:            req.Teamwork,
			Initiative:          req.Initiative,
			Notes:               req.Notes,
			SubmittedAt:         time.Now(),
		}
	} else {
		penilaian.AttendanceScoreAuto = attendanceScore
		penilaian.Discipline = req.Discipline
		penilaian.Responsibility = req.Responsibility
		penilaian.Teamwork = req.Teamwork
		penilaian.Initiative = req.Initiative
		penilaian.Notes = req.Notes
		penilaian.SubmittedAt = time.Now()
	}

	penilaian.FinalScore, penilaian.FinalGrade = calculateFinalScore(penilaian)

	if result.Error != nil {
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
	query := database.DB.Preload("Student")

	switch role {
	case "dudi":
		userID, _ := c.Get("user_id")
		uid, _ := uuid.Parse(userID.(string))
		query = query.Where("dudi_id = ?", uid)
	case "teacher", "admin", "admin_jurusan":
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

func calculateAttendanceScore(studentID uuid.UUID) float64 {
	var count int64
	database.DB.Model(&models.Absensi{}).
		Where("student_id = ? AND status IN ('hadir','terlambat')", studentID).
		Count(&count)

	var total int64
	database.DB.Model(&models.Absensi{}).Where("student_id = ?", studentID).Count(&total)

	if total == 0 {
		return 0
	}

	return float64(count) / float64(total) * 100
}

func calculateFinalScore(p models.Penilaian) (float64, string) {
	attendanceWeight := 0.3
	manualWeight := 0.7

	manualAvg := float64(p.Discipline+p.Responsibility+p.Teamwork+p.Initiative) / 4.0
	score := (p.AttendanceScoreAuto * attendanceWeight) + (manualAvg/5*100)*manualWeight

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
