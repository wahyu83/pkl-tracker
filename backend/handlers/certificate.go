package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type CertificateHandler struct{}

func NewCertificateHandler() *CertificateHandler {
	return &CertificateHandler{}
}

type CertificateStudent struct {
	ID        uuid.UUID `json:"id"`
	FullName  string    `json:"full_name"`
	NisNipNik string    `json:"nis_nip_nik"`
	Jurusan   string    `json:"jurusan"`
}

type CertificateDudi struct {
	CompanyName string `json:"company_name"`
	Address     string `json:"address"`
	PicName     string `json:"pic_name"`
}

type CertificateScore struct {
	FinalScore          float64    `json:"final_score"`
	FinalGrade          string     `json:"final_grade"`
	AttendanceScoreAuto float64    `json:"attendance_score_auto"`
	Notes               string     `json:"notes"`
	SubmittedAt         *time.Time `json:"submitted_at"`
	Assessed            bool       `json:"assessed"`
}

type CertificateAbsensi struct {
	Hadir     int `json:"hadir"`
	Terlambat int `json:"terlambat"`
	Izin      int `json:"izin"`
	Sakit     int `json:"sakit"`
	Total     int `json:"total"`
}

type Certificate struct {
	CertificateNumber string             `json:"certificate_number"`
	Student           CertificateStudent `json:"student"`
	Dudi              *CertificateDudi   `json:"dudi"`
	Score             CertificateScore   `json:"score"`
	Absensi           CertificateAbsensi `json:"absensi"`
}

// resolvePeriode mengambil periode berdasarkan id, atau periode aktif bila kosong.
func resolvePeriode(periodeID string) *models.Periode {
	if periodeID != "" {
		id, err := uuid.Parse(periodeID)
		if err != nil {
			return nil
		}
		var p models.Periode
		if database.DB.First(&p, "id = ?", id).Error != nil {
			return nil
		}
		return &p
	}

	var active models.Periode
	if database.DB.Where("is_active = ?", true).First(&active).Error != nil {
		return nil
	}
	return &active
}

// periodeStartYear mengambil tahun awal dari "2025/2026" -> "2025".
func periodeStartYear(tahunPelajaran string) string {
	if idx := strings.Index(tahunPelajaran, "/"); idx > 0 {
		return tahunPelajaran[:idx]
	}
	if tahunPelajaran != "" {
		return tahunPelajaran
	}
	return fmt.Sprintf("%d", time.Now().Year())
}

func buildCertificateNumber(prefix, tahunPelajaran, nis string) string {
	if prefix == "" {
		prefix = "421.5/PKL"
	}
	return fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(prefix, "/"), periodeStartYear(tahunPelajaran), nis)
}

// attendanceSummary menghitung rekap absensi per siswa pada rentang periode.
type attendanceRow struct {
	StudentID uuid.UUID
	Status    string
	Total     int
}

func attendanceSummary(studentIDs []uuid.UUID, periode *models.Periode) map[uuid.UUID]CertificateAbsensi {
	result := make(map[uuid.UUID]CertificateAbsensi)
	if len(studentIDs) == 0 {
		return result
	}

	query := database.DB.Table("absensis").
		Select("student_id, status, count(*) as total").
		Where("student_id IN ?", studentIDs)

	if periode != nil {
		query = query.Where("timestamp >= ? AND timestamp < ?::date + interval '1 day'",
			periode.StartDate, periode.EndDate)
	}

	var rows []attendanceRow
	query.Group("student_id, status").Scan(&rows)

	for _, r := range rows {
		entry := result[r.StudentID]
		switch r.Status {
		case "hadir":
			entry.Hadir = r.Total
		case "terlambat":
			entry.Terlambat = r.Total
		case "izin":
			entry.Izin = r.Total
		case "sakit":
			entry.Sakit = r.Total
		}
		entry.Total += r.Total
		result[r.StudentID] = entry
	}

	return result
}

func penilaianByStudent(studentIDs []uuid.UUID) map[uuid.UUID]models.Penilaian {
	result := make(map[uuid.UUID]models.Penilaian)
	if len(studentIDs) == 0 {
		return result
	}

	var list []models.Penilaian
	database.DB.Where("student_id IN ?", studentIDs).Find(&list)
	for _, p := range list {
		result[p.StudentID] = p
	}
	return result
}

func buildCertificate(student models.User, periode *models.Periode, prefix string,
	penilaian *models.Penilaian, absensi CertificateAbsensi) Certificate {

	cert := Certificate{
		CertificateNumber: buildCertificateNumber(prefix, periodeTahun(periode), student.NisNipNik),
		Student: CertificateStudent{
			ID:        student.ID,
			FullName:  student.FullName,
			NisNipNik: student.NisNipNik,
			Jurusan:   student.Jurusan,
		},
		Absensi: absensi,
	}

	if student.DUDI != nil {
		cert.Dudi = &CertificateDudi{
			CompanyName: student.DUDI.CompanyName,
			Address:     student.DUDI.Address,
			PicName:     student.DUDI.PicName,
		}
	}

	if penilaian != nil {
		submitted := penilaian.SubmittedAt
		cert.Score = CertificateScore{
			FinalScore:          penilaian.FinalScore,
			FinalGrade:          penilaian.FinalGrade,
			AttendanceScoreAuto: penilaian.AttendanceScoreAuto,
			Notes:               penilaian.Notes,
			SubmittedAt:         &submitted,
			Assessed:            true,
		}
	}

	return cert
}

func periodeTahun(periode *models.Periode) string {
	if periode == nil {
		return ""
	}
	return periode.TahunPelajaran
}

func (h *CertificateHandler) GetOne(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	jurusan, _ := c.Get("jurusan")
	uid, _ := uuid.Parse(userID.(string))

	studentID, err := uuid.Parse(c.Param("studentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID siswa tidak valid"})
		return
	}

	query := database.DB.Preload("DUDI").Where("role = 'student'")

	switch role {
	case "admin":
		// semua siswa
	case "admin_jurusan":
		if jurusan == nil || jurusan.(string) == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		query = query.Where("jurusan = ?", jurusan.(string))
	case "teacher":
		query = query.Where("teacher_id = ?", uid)
	case "dudi":
		var dudiUser models.User
		if database.DB.First(&dudiUser, "id = ?", uid).Error != nil || dudiUser.DudiID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Profil DUDI tidak ditemukan"})
			return
		}
		query = query.Where("dudi_id = ?", *dudiUser.DudiID)
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var student models.User
	if err := query.Where("id = ?", studentID).First(&student).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Siswa tidak ditemukan atau bukan dalam akses Anda"})
		return
	}

	periode := resolvePeriode(c.Query("periode_id"))
	profile := getOrCreateSchoolProfile()

	penilaianMap := penilaianByStudent([]uuid.UUID{student.ID})
	absensiMap := attendanceSummary([]uuid.UUID{student.ID}, periode)

	var penilaian *models.Penilaian
	if p, ok := penilaianMap[student.ID]; ok {
		penilaian = &p
	}

	cert := buildCertificate(student, periode, profile.CertificatePrefix, penilaian, absensiMap[student.ID])

	c.JSON(http.StatusOK, gin.H{
		"school":  profile,
		"periode": periode,
		"data":    cert,
	})
}

func (h *CertificateHandler) List(c *gin.Context) {
	role, _ := c.Get("role")
	userID, _ := c.Get("user_id")
	jurusan, _ := c.Get("jurusan")
	uid, _ := uuid.Parse(userID.(string))

	query := database.DB.Preload("DUDI").Where("role = 'student'")

	switch role {
	case "admin":
		// semua siswa
	case "admin_jurusan":
		if jurusan == nil || jurusan.(string) == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
		query = query.Where("jurusan = ?", jurusan.(string))
	case "teacher":
		query = query.Where("teacher_id = ?", uid)
	case "dudi":
		var dudiUser models.User
		if database.DB.First(&dudiUser, "id = ?", uid).Error != nil || dudiUser.DudiID == nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "Profil DUDI tidak ditemukan"})
			return
		}
		query = query.Where("dudi_id = ?", *dudiUser.DudiID)
	default:
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	// Filter opsional
	if dudiFilter := c.Query("dudi_id"); dudiFilter != "" {
		if id, err := uuid.Parse(dudiFilter); err == nil {
			query = query.Where("dudi_id = ?", id)
		}
	}
	if jurusanFilter := c.Query("jurusan"); jurusanFilter != "" {
		query = query.Where("jurusan = ?", jurusanFilter)
	}
	if search := c.Query("search"); search != "" {
		q := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR nis_nip_nik ILIKE ?", q, q)
	}

	var students []models.User
	query.Order("full_name ASC").Find(&students)

	ids := make([]uuid.UUID, 0, len(students))
	for _, s := range students {
		ids = append(ids, s.ID)
	}

	periode := resolvePeriode(c.Query("periode_id"))
	profile := getOrCreateSchoolProfile()
	penilaianMap := penilaianByStudent(ids)
	absensiMap := attendanceSummary(ids, periode)

	certificates := make([]Certificate, 0, len(students))
	for _, s := range students {
		var penilaian *models.Penilaian
		if p, ok := penilaianMap[s.ID]; ok {
			penilaian = &p
		}
		certificates = append(certificates, buildCertificate(s, periode, profile.CertificatePrefix, penilaian, absensiMap[s.ID]))
	}

	c.JSON(http.StatusOK, gin.H{
		"school":  profile,
		"periode": periode,
		"total":   len(certificates),
		"data":    certificates,
	})
}
