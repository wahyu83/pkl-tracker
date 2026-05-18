package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

func checkAdminAccess(c *gin.Context) (isAdmin bool, adminJurusan string) {
	role, _ := c.Get("role")
	jurusan, _ := c.Get("jurusan")
	if role == "admin" {
		return true, ""
	}
	if role == "admin_jurusan" && jurusan != "" {
		return false, jurusan.(string)
	}
	return false, ""
}

// --- Users ---

type CreateUserRequest struct {
	FullName  string `json:"full_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
	NisNipNik string `json:"nis_nip_nik" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	Role      string `json:"role" binding:"required,oneof=student teacher dudi admin admin_jurusan"`
	Jurusan   string `json:"jurusan"`
	DudiID    string `json:"dudi_id"`
}

type UpdateUserRequest struct {
	FullName  string `json:"full_name"`
	Email     string `json:"email"`
	NisNipNik string `json:"nis_nip_nik"`
	Password  string `json:"password"`
	Role      string `json:"role" binding:"omitempty,oneof=student teacher dudi admin admin_jurusan"`
	Jurusan   string `json:"jurusan"`
	DudiID    string `json:"dudi_id"`
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var users []models.User
	query := database.DB.Model(&models.User{})

	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}
	if jurusanFilter := c.Query("jurusan"); jurusanFilter != "" {
		query = query.Where("jurusan = ?", jurusanFilter)
	}
	if roleFilter := c.Query("role"); roleFilter != "" {
		query = query.Where("role = ?", roleFilter)
	}
	if search := c.Query("search"); search != "" {
		q := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR email ILIKE ? OR nis_nip_nik ILIKE ?", q, q, q)
	}

	query.Order("created_at DESC").Find(&users)
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *AdminHandler) CreateUser(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.User
	if database.DB.Where("email = ? OR nis_nip_nik = ?", req.Email, req.NisNipNik).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email or NIS/NIP/NIK already in use"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	jurusan := req.Jurusan
	if !isAdmin && adminJurusan != "" {
		jurusan = adminJurusan
	}

	user := models.User{
		FullName:     req.FullName,
		Email:        req.Email,
		NisNipNik:    req.NisNipNik,
		PasswordHash: string(hash),
		Role:         req.Role,
		Jurusan:      jurusan,
	}

	if req.DudiID != "" {
		dudiID, err := uuid.Parse(req.DudiID)
		if err == nil {
			user.DudiID = &dudiID
		}
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created", "data": user})
}

func (h *AdminHandler) UpdateUser(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	if database.DB.First(&user, "id = ?", id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if !isAdmin && adminJurusan != "" && user.Jurusan != adminJurusan {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot edit user outside your jurusan"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Email != "" && req.Email != user.Email {
		var dup models.User
		if database.DB.Where("email = ? AND id != ?", req.Email, id).First(&dup).Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already in use"})
			return
		}
		user.Email = req.Email
	}
	if req.NisNipNik != "" && req.NisNipNik != user.NisNipNik {
		var dup models.User
		if database.DB.Where("nis_nip_nik = ? AND id != ?", req.NisNipNik, id).First(&dup).Error == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "NIS/NIP/NIK already in use"})
			return
		}
		user.NisNipNik = req.NisNipNik
	}
	if req.Password != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		user.PasswordHash = string(hash)
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if isAdmin && req.Jurusan != "" {
		user.Jurusan = req.Jurusan
	}
	if req.DudiID != "" {
		dudiID, err := uuid.Parse(req.DudiID)
		if err == nil {
			user.DudiID = &dudiID
		}
	} else if req.DudiID == "" && c.Request.Body != nil {
		user.DudiID = nil
	}

	database.DB.Save(&user)
	c.JSON(http.StatusOK, gin.H{"message": "User updated", "data": user})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	query := database.DB.Where("id = ?", id)
	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}

	result := query.Delete(&models.User{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user: " + result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

type BulkDeleteRequest struct {
	IDs []string `json:"ids" binding:"required"`
}

func (h *AdminHandler) BulkDeleteUsers(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req BulkDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No IDs provided"})
		return
	}

	query := database.DB.Where("id IN ?", req.IDs)
	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}

	result := query.Delete(&models.User{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus user: " + result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Users deleted", "deleted": result.RowsAffected})
}

// --- DUDI ---

type DUDIRequest struct {
	CompanyName   string  `json:"company_name" binding:"required"`
	Address       string  `json:"address"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	RadiusAllowed int     `json:"radius_allowed"`
	PicName       string  `json:"pic_name"`
	Phone         string  `json:"phone"`
	Jurusan       string  `json:"jurusan"`
}

func (h *AdminHandler) ListDUDI(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var dudiList []models.DUDI
	query := database.DB.Model(&models.DUDI{})

	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}
	if jurusanFilter := c.Query("jurusan"); jurusanFilter != "" {
		query = query.Where("jurusan = ?", jurusanFilter)
	}
	if search := c.Query("search"); search != "" {
		q := "%" + search + "%"
		query = query.Where("company_name ILIKE ? OR pic_name ILIKE ?", q, q)
	}

	query.Order("created_at DESC").Find(&dudiList)

	type DUDIWithStudentCount struct {
		models.DUDI
		StudentCount int64 `json:"student_count"`
	}

	result := make([]DUDIWithStudentCount, len(dudiList))
	for i, d := range dudiList {
		var count int64
		database.DB.Model(&models.User{}).Where("dudi_id = ?", d.ID).Count(&count)
		result[i] = DUDIWithStudentCount{DUDI: d, StudentCount: count}
	}

	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *AdminHandler) CreateDUDI(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var req DUDIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jurusan := req.Jurusan
	if !isAdmin && adminJurusan != "" {
		jurusan = adminJurusan
	}

	dudi := models.DUDI{
		CompanyName:   req.CompanyName,
		Address:       req.Address,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		RadiusAllowed: req.RadiusAllowed,
		PicName:       req.PicName,
		Phone:         req.Phone,
		Jurusan:       jurusan,
	}

	if dudi.RadiusAllowed == 0 {
		dudi.RadiusAllowed = 500
	}

	if err := database.DB.Create(&dudi).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create DUDI"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "DUDI created", "data": dudi})
}

func (h *AdminHandler) UpdateDUDI(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DUDI ID"})
		return
	}

	query := database.DB.Where("id = ?", id)
	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}

	var dudi models.DUDI
	if query.First(&dudi).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "DUDI not found"})
		return
	}

	var req DUDIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dudi.CompanyName = req.CompanyName
	dudi.Address = req.Address
	dudi.Latitude = req.Latitude
	dudi.Longitude = req.Longitude
	dudi.RadiusAllowed = req.RadiusAllowed
	dudi.PicName = req.PicName
	dudi.Phone = req.Phone
	if isAdmin && req.Jurusan != "" {
		dudi.Jurusan = req.Jurusan
	}

	database.DB.Save(&dudi)
	c.JSON(http.StatusOK, gin.H{"message": "DUDI updated", "data": dudi})
}

func (h *AdminHandler) DeleteDUDI(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DUDI ID"})
		return
	}

	query := database.DB.Where("id = ?", id)
	if !isAdmin && adminJurusan != "" {
		query = query.Where("jurusan = ?", adminJurusan)
	}

	if query.Delete(&models.DUDI{}).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "DUDI not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "DUDI deleted"})
}

// --- Periode ---

type PeriodeRequest struct {
	TahunPelajaran string `json:"tahun_pelajaran" binding:"required"`
	Semester       string `json:"semester" binding:"required,oneof=ganjil genap"`
	StartDate      string `json:"start_date" binding:"required"`
	EndDate        string `json:"end_date" binding:"required"`
}

func (h *AdminHandler) ListPeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		role, _ := c.Get("role")
		if role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}
	_ = adminJurusan

	var periods []models.Periode
	database.DB.Order("start_date DESC").Find(&periods)
	if periods == nil {
		periods = []models.Periode{}
	}
	c.JSON(http.StatusOK, gin.H{"data": periods})
}

func (h *AdminHandler) ActivePeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		role, _ := c.Get("role")
		if role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var periode models.Periode
	if database.DB.Where("is_active = ?", true).First(&periode).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No active period"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": periode})
}

func (h *AdminHandler) CreatePeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	_ = adminJurusan

	var req PeriodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format (use YYYY-MM-DD)"})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format (use YYYY-MM-DD)"})
		return
	}

	periode := models.Periode{
		TahunPelajaran: req.TahunPelajaran,
		Semester:       req.Semester,
		StartDate:      startDate,
		EndDate:        endDate,
	}

	if err := database.DB.Create(&periode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create period"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Period created", "data": periode})
}

func (h *AdminHandler) UpdatePeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	_ = adminJurusan

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	var periode models.Periode
	if database.DB.First(&periode, "id = ?", id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Period not found"})
		return
	}

	var req PeriodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}
	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	periode.TahunPelajaran = req.TahunPelajaran
	periode.Semester = req.Semester
	periode.StartDate = startDate
	periode.EndDate = endDate

	database.DB.Save(&periode)
	c.JSON(http.StatusOK, gin.H{"message": "Period updated", "data": periode})
}

func (h *AdminHandler) ActivatePeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	_ = adminJurusan

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	database.DB.Model(&models.Periode{}).Where("is_active = ?", true).Update("is_active", false)

	database.DB.Model(&models.Periode{}).Where("id = ?", id).Update("is_active", true)

	var periode models.Periode
	database.DB.First(&periode, "id = ?", id)

	c.JSON(http.StatusOK, gin.H{"message": "Period activated", "data": periode})
}

func (h *AdminHandler) DeletePeriode(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	_ = adminJurusan

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid period ID"})
		return
	}

	if database.DB.Delete(&models.Periode{}, "id = ?", id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Period not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Period deleted"})
}
