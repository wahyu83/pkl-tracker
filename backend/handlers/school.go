package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

// timestampSuffix dipakai sebagai query cache-buster pada URL logo.
func timestampSuffix() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

type SchoolHandler struct{}

func NewSchoolHandler() *SchoolHandler {
	return &SchoolHandler{}
}

// defaultSchoolProfile dipakai saat pengaturan belum pernah diisi.
func defaultSchoolProfile() models.SchoolProfile {
	return models.SchoolProfile{
		Name:              "SMKN 1 Arahan",
		City:              "Indramayu",
		CertificatePrefix: "421.5/PKL",
	}
}

func getOrCreateSchoolProfile() models.SchoolProfile {
	var profile models.SchoolProfile
	if database.DB.First(&profile).Error == nil {
		return profile
	}

	profile = defaultSchoolProfile()
	database.DB.Create(&profile)
	return profile
}

func (h *SchoolHandler) Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": getOrCreateSchoolProfile()})
}

type SchoolProfileRequest struct {
	Name              string `json:"name"`
	Address           string `json:"address"`
	City              string `json:"city"`
	HeadmasterName    string `json:"headmaster_name"`
	HeadmasterNip     string `json:"headmaster_nip"`
	CertificatePrefix string `json:"certificate_prefix"`
}

func (h *SchoolHandler) Update(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengubah pengaturan sekolah"})
		return
	}

	var req SchoolProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Data pengaturan tidak valid"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama sekolah wajib diisi"})
		return
	}

	profile := getOrCreateSchoolProfile()
	profile.Name = strings.TrimSpace(req.Name)
	profile.Address = strings.TrimSpace(req.Address)
	profile.City = strings.TrimSpace(req.City)
	profile.HeadmasterName = strings.TrimSpace(req.HeadmasterName)
	profile.HeadmasterNip = strings.TrimSpace(req.HeadmasterNip)
	profile.CertificatePrefix = strings.TrimSpace(req.CertificatePrefix)

	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengaturan sekolah disimpan", "data": profile})
}

var allowedLogoExt = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

func (h *SchoolHandler) UploadLogo(c *gin.Context) {
	if role, _ := c.Get("role"); role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Hanya admin yang dapat mengubah logo sekolah"})
		return
	}

	// Batasi ukuran body (3 MB) untuk menghindari unggahan besar.
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 3<<20)

	file, header, err := c.Request.FormFile("logo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File logo tidak ditemukan"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedLogoExt[ext] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Format logo harus PNG, JPG, atau WEBP"})
		return
	}

	dir := filepath.Join("uploads", "school")
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat folder upload"})
		return
	}

	// Nama tetap supaya logo lama tertimpa dan tidak menumpuk.
	filename := "logo" + ext
	dest := filepath.Join(dir, filename)

	out, err := os.Create(dest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan logo"})
		return
	}
	defer out.Close()

	if _, err := out.ReadFrom(file); err != nil {
		os.Remove(dest)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan logo"})
		return
	}

	profile := getOrCreateSchoolProfile()

	// Hapus logo lama dengan ekstensi berbeda agar tidak tersisa.
	for e := range allowedLogoExt {
		if e != ext {
			old := filepath.Join(dir, "logo"+e)
			if old != dest {
				os.Remove(old)
			}
		}
	}

	profile.LogoURL = "/uploads/school/" + filename + "?v=" + timestampSuffix()
	if err := database.DB.Save(&profile).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pengaturan logo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Logo sekolah disimpan", "data": profile})
}
