package handlers

import (
	"fmt"
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
	if isAdmin {
		user.Jurusan = req.Jurusan
	} else if adminJurusan != "" {
		user.Jurusan = adminJurusan
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

// --- Dashboard ---

func (h *AdminHandler) Dashboard(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	type RoleCount struct {
		Role  string `json:"role"`
		Count int64  `json:"count"`
	}
	var roleCounts []RoleCount
	userQuery := database.DB.Model(&models.User{})
	if !isAdmin && adminJurusan != "" {
		userQuery = userQuery.Where("jurusan = ?", adminJurusan)
	}
	userQuery.Select("role, count(*) as count").Group("role").Scan(&roleCounts)

	totalSiswa := int64(0)
	totalGuru := int64(0)
	totalAdmin := int64(0)
	for _, rc := range roleCounts {
		switch rc.Role {
		case "student":
			totalSiswa = rc.Count
		case "teacher":
			totalGuru = rc.Count
		case "admin", "admin_jurusan":
			totalAdmin += rc.Count
		}
	}

	var totalDudi int64
	dudiQuery := database.DB.Model(&models.DUDI{})
	if !isAdmin && adminJurusan != "" {
		dudiQuery = dudiQuery.Where("jurusan = ?", adminJurusan)
	}
	dudiQuery.Count(&totalDudi)

	var activePeriodCount int64
	database.DB.Model(&models.Periode{}).Where("is_active = ?", true).Count(&activePeriodCount)

	type Activity struct {
		Type string `json:"type"`
		Text string `json:"text"`
		Time string `json:"time"`
	}
	activities := []Activity{}

	type AbsensiActivity struct {
		FullName  string
		CreatedAt time.Time
	}
	var recentAbsensi []AbsensiActivity
	absQuery := database.DB.Table("absensis").
		Select("users.full_name, absensis.created_at").
		Joins("JOIN users ON users.id = absensis.student_id")
	if !isAdmin && adminJurusan != "" {
		absQuery = absQuery.Where("users.jurusan = ?", adminJurusan)
	}
	absQuery.Order("absensis.created_at DESC").Limit(3).Scan(&recentAbsensi)
	for _, a := range recentAbsensi {
		activities = append(activities, Activity{
			Type: "absensi",
			Text: "Siswa \"" + a.FullName + "\" melakukan absensi",
			Time: timeAgo(a.CreatedAt),
		})
	}

	type JurnalActivity struct {
		FullName  string
		CreatedAt time.Time
	}
	var recentJurnal []JurnalActivity
	jrnQuery := database.DB.Table("jurnals").
		Select("users.full_name, jurnals.created_at").
		Joins("JOIN users ON users.id = jurnals.student_id")
	if !isAdmin && adminJurusan != "" {
		jrnQuery = jrnQuery.Where("users.jurusan = ?", adminJurusan)
	}
	jrnQuery.Order("jurnals.created_at DESC").Limit(3).Scan(&recentJurnal)
	for _, j := range recentJurnal {
		activities = append(activities, Activity{
			Type: "jurnal",
			Text: "Siswa \"" + j.FullName + "\" menulis jurnal baru",
			Time: timeAgo(j.CreatedAt),
		})
	}

	type PenilaianActivity struct {
		FullName    string
		SubmittedAt time.Time
	}
	var recentNilai []PenilaianActivity
	nilQuery := database.DB.Table("penilaians").
		Select("users.full_name, penilaians.submitted_at").
		Joins("JOIN users ON users.id = penilaians.student_id")
	if !isAdmin && adminJurusan != "" {
		nilQuery = nilQuery.Where("users.jurusan = ?", adminJurusan)
	}
	nilQuery.Where("penilaians.submitted_at IS NOT NULL").Order("penilaians.submitted_at DESC").Limit(3).Scan(&recentNilai)
	for _, n := range recentNilai {
		activities = append(activities, Activity{
			Type: "penilaian",
			Text: "Siswa \"" + n.FullName + "\" mendapat nilai PKL",
			Time: timeAgo(n.SubmittedAt),
		})
	}

	distributions := []map[string]interface{}{
		{"label": "Siswa", "count": totalSiswa, "role": "student"},
		{"label": "Guru", "count": totalGuru, "role": "teacher"},
		{"label": "DUDI", "count": totalDudi, "role": "dudi"},
		{"label": "Admin", "count": totalAdmin, "role": "admin"},
	}
	totalUsers := totalSiswa + totalGuru + totalDudi + totalAdmin
	for i, d := range distributions {
		if totalUsers > 0 {
			distributions[i]["percent"] = int(float64(d["count"].(int64)) / float64(totalUsers) * 100)
		} else {
			distributions[i]["percent"] = 0
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"stats": map[string]interface{}{
			"total_siswa":    totalSiswa,
			"total_guru":     totalGuru,
			"total_dudi":     totalDudi,
			"total_admin":    totalAdmin,
			"active_period":  activePeriodCount,
		},
		"distributions":    distributions,
		"recent_activities": activities,
	})
}

func timeAgo(t time.Time) string {
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "Baru saja"
	case d < time.Hour:
		m := int(d.Minutes())
		if m == 1 {
			return "1 menit lalu"
		}
		return fmt.Sprintf("%d menit lalu", m)
	case d < 24*time.Hour:
		h := int(d.Hours())
		if h == 1 {
			return "1 jam lalu"
		}
		return fmt.Sprintf("%d jam lalu", h)
	default:
		days := int(d.Hours() / 24)
		if days == 1 {
			return "1 hari lalu"
		}
		return fmt.Sprintf("%d hari lalu", days)
	}
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
	DudiNIK       string  `json:"dudi_nik"`
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
		StudentCount int64    `json:"student_count"`
		DudiNIKs     []string `json:"dudi_niks"`
	}

	result := make([]DUDIWithStudentCount, len(dudiList))
	for i, d := range dudiList {
		var count int64
		database.DB.Model(&models.User{}).Where("dudi_id = ?", d.ID).Count(&count)
		result[i] = DUDIWithStudentCount{DUDI: d, StudentCount: count}

		var dudiUsers []models.User
		database.DB.Where("dudi_id = ? AND role = 'dudi'", d.ID).Find(&dudiUsers)
		niks := make([]string, len(dudiUsers))
		for j, u := range dudiUsers {
			niks[j] = u.NisNipNik
		}
		result[i].DudiNIKs = niks
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

	if req.DudiNIK != "" {
		linkOrCreateDudiUser(dudi.ID, req.DudiNIK, req.Jurusan, req.CompanyName)
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
	if isAdmin {
		dudi.Jurusan = req.Jurusan
	} else if adminJurusan != "" {
		dudi.Jurusan = adminJurusan
	}

	database.DB.Save(&dudi)

	if req.DudiNIK != "" {
		linkOrCreateDudiUser(dudi.ID, req.DudiNIK, dudi.Jurusan, dudi.CompanyName)
	}

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

// --- Jurusan ---

type JurusanRequest struct {
	Nama string `json:"nama" binding:"required"`
	Kode string `json:"kode" binding:"required"`
}

func (h *AdminHandler) ListJurusan(c *gin.Context) {
	_, adminJurusan := checkAdminAccess(c)
	role, _ := c.Get("role")
	if role != "admin" && role != "admin_jurusan" && role != "teacher" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}
	_ = adminJurusan

	var list []models.Jurusan
	database.DB.Order("nama ASC").Find(&list)
	if list == nil {
		list = []models.Jurusan{}
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

func (h *AdminHandler) CreateJurusan(c *gin.Context) {
	isAdmin, _ := checkAdminAccess(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can manage jurusan"})
		return
	}

	var req JurusanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existing models.Jurusan
	if database.DB.Where("kode = ? OR nama = ?", req.Kode, req.Nama).First(&existing).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Jurusan dengan kode atau nama tersebut sudah ada"})
		return
	}

	j := models.Jurusan{Nama: req.Nama, Kode: req.Kode}
	if err := database.DB.Create(&j).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat jurusan"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Jurusan created", "data": j})
}

func (h *AdminHandler) UpdateJurusan(c *gin.Context) {
	isAdmin, _ := checkAdminAccess(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can manage jurusan"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var j models.Jurusan
	if database.DB.First(&j, "id = ?", id).Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurusan not found"})
		return
	}

	var req JurusanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var dup models.Jurusan
	if database.DB.Where("(kode = ? OR nama = ?) AND id != ?", req.Kode, req.Nama, id).First(&dup).Error == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Jurusan dengan kode atau nama tersebut sudah ada"})
		return
	}

	j.Nama = req.Nama
	j.Kode = req.Kode
	database.DB.Save(&j)
	c.JSON(http.StatusOK, gin.H{"message": "Jurusan updated", "data": j})
}

func (h *AdminHandler) DeleteJurusan(c *gin.Context) {
	isAdmin, _ := checkAdminAccess(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can manage jurusan"})
		return
	}

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("jurusan = (SELECT kode FROM jurusans WHERE id = ?)", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "Tidak bisa menghapus jurusan yang masih digunakan oleh pengguna"})
		return
	}

	if database.DB.Delete(&models.Jurusan{}, "id = ?", id).RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Jurusan not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jurusan deleted"})
}

func linkOrCreateDudiUser(dudiID uuid.UUID, nik, jurusan, companyName string) {
	var existing models.User
	if database.DB.Where("nis_nip_nik = ? AND role = 'dudi'", nik).First(&existing).Error == nil {
		existing.DudiID = &dudiID
		if existing.Jurusan == "" && jurusan != "" {
			existing.Jurusan = jurusan
		}
		database.DB.Save(&existing)
		return
	}

	defaultPass, _ := bcrypt.GenerateFromPassword([]byte("pkl123456"), bcrypt.DefaultCost)
	dudiUser := models.User{
		FullName:     companyName,
		Email:        nik + "@dudi.local",
		PasswordHash: string(defaultPass),
		Role:         "dudi",
		NisNipNik:    nik,
		Jurusan:      jurusan,
		DudiID:       &dudiID,
	}
	database.DB.Where("email = ?", dudiUser.Email).FirstOrCreate(&dudiUser)
}
