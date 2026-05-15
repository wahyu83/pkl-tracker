package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type JurnalHandler struct{}

func NewJurnalHandler() *JurnalHandler {
	return &JurnalHandler{}
}

type JurnalRequest struct {
	Date       string `json:"date" binding:"required"`
	Activity   string `json:"activity" binding:"required"`
	Reflection string `json:"reflection"`
}

type CommentRequest struct {
	JurnalID string `json:"jurnal_id" binding:"required"`
	Comment  string `json:"comment" binding:"required"`
}

func (h *JurnalHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can create journals"})
		return
	}

	var req JurnalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Activity == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Activity is required"})
		return
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		date = time.Now()
	}

	uid, _ := uuid.Parse(userID.(string))

	jurnal := models.Jurnal{
		StudentID:  uid,
		Date:       date,
		Activity:   req.Activity,
		Reflection: req.Reflection,
	}

	if err := database.DB.Create(&jurnal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save journal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Journal created", "data": jurnal})
}

func (h *JurnalHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	uid, _ := uuid.Parse(userID.(string))

	var jurnalList []models.Jurnal
	query := database.DB.Preload("Student").Order("date DESC")

	if role == "student" {
		query = query.Where("student_id = ?", uid)
	} else {
		studentID := c.Query("student_id")
		if studentID != "" {
			query = query.Where("student_id = ?", studentID)
		}
	}

	query.Find(&jurnalList)

	c.JSON(http.StatusOK, gin.H{"data": jurnalList})
}

func (h *JurnalHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	jurnalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var jurnal models.Jurnal
	if err := database.DB.Preload("Student").First(&jurnal, "id = ?", jurnalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": jurnal})
}

func (h *JurnalHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	if role != "student" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only students can edit journals"})
		return
	}

	id := c.Param("id")
	jurnalID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var jurnal models.Jurnal
	if err := database.DB.First(&jurnal, "id = ?", jurnalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	uid, _ := uuid.Parse(userID.(string))
	if jurnal.StudentID != uid {
		c.JSON(http.StatusForbidden, gin.H{"error": "Can only edit your own journals"})
		return
	}

	if jurnal.TeacherComment != "" || jurnal.DudiComment != "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot edit journal after comments are made"})
		return
	}

	var req JurnalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	date, _ := time.Parse("2006-01-02", req.Date)
	jurnal.Date = date
	jurnal.Activity = req.Activity
	jurnal.Reflection = req.Reflection

	database.DB.Save(&jurnal)
	c.JSON(http.StatusOK, gin.H{"message": "Journal updated", "data": jurnal})
}

func (h *JurnalHandler) Comment(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jurnalID, err := uuid.Parse(req.JurnalID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jurnal ID"})
		return
	}

	var jurnal models.Jurnal
	if err := database.DB.First(&jurnal, "id = ?", jurnalID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	switch role {
	case "teacher":
		jurnal.TeacherComment = req.Comment
	case "dudi":
		var student models.User
		if database.DB.First(&student, "id = ?", jurnal.StudentID).Error == nil {
			uid, _ := uuid.Parse(userID.(string))
			var dudiUser models.User
			if database.DB.First(&dudiUser, "id = ?", uid).Error == nil {
				if dudiUser.DudiID == nil || *dudiUser.DudiID != *student.DudiID {
					c.JSON(http.StatusForbidden, gin.H{"error": "You can only comment on students at your company"})
					return
				}
			}
		}
		jurnal.DudiComment = req.Comment
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Only teachers and DUDI can comment"})
		return
	}

	database.DB.Save(&jurnal)
	c.JSON(http.StatusOK, gin.H{"message": "Comment added", "data": jurnal})
}
