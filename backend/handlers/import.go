package handlers

import (
	"encoding/csv"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"pkl-tracker/database"
	"pkl-tracker/models"
)

type ImportHandler struct{}

func NewImportHandler() *ImportHandler {
	return &ImportHandler{}
}

func (h *ImportHandler) ImportSiswa(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can import"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV empty or invalid"})
		return
	}

	colIndex := mapCSVHeaders(headers)
	required := []string{"full_name", "email", "nis"}
	for _, r := range required {
		if _, ok := colIndex[r]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom '" + r + "' wajib ada di CSV"})
			return
		}
	}

	imported := 0
	skipped := 0
	errors := []string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			errors = append(errors, "Gagal membaca baris CSV")
			continue
		}

		fullName := getCol(record, colIndex, "full_name")
		email := getCol(record, colIndex, "email")
		nis := getCol(record, colIndex, "nis")
		password := getCol(record, colIndex, "password")
		dudiNIK := getCol(record, colIndex, "dudi_nik")

		if fullName == "" || email == "" || nis == "" {
			skipped++
			continue
		}

		if password == "" {
			password = "pkl123456"
		}

		var existing models.User
		if database.DB.Where("email = ? OR nis_nip_nik = ?", email, nis).First(&existing).Error == nil {
			skipped++
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			skipped++
			continue
		}

		user := models.User{
			FullName:     fullName,
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "student",
			NisNipNik:    nis,
		}

		if dudiNIK != "" {
			var dudi models.DUDI
			var dudiUser models.User
			if database.DB.Where("nis_nip_nik = ? AND role = 'dudi'", dudiNIK).First(&dudiUser).Error == nil {
				if dudiUser.DudiID != nil {
					user.DudiID = dudiUser.DudiID
				}
			} else if database.DB.Where("pic_name = ?", dudiNIK).First(&dudi).Error == nil {
				user.DudiID = &dudi.ID
			}
		}

		if err := database.DB.Create(&user).Error; err != nil {
			skipped++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import selesai",
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
	})
}

func (h *ImportHandler) ImportGuru(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can import"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV empty or invalid"})
		return
	}

	colIndex := mapCSVHeaders(headers)
	required := []string{"full_name", "email", "nip"}
	for _, r := range required {
		if _, ok := colIndex[r]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom '" + r + "' wajib ada di CSV"})
			return
		}
	}

	imported := 0
	skipped := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		fullName := getCol(record, colIndex, "full_name")
		email := getCol(record, colIndex, "email")
		nip := getCol(record, colIndex, "nip")
		password := getCol(record, colIndex, "password")

		if fullName == "" || email == "" || nip == "" {
			skipped++
			continue
		}

		if password == "" {
			password = "pkl123456"
		}

		var existing models.User
		if database.DB.Where("email = ? OR nis_nip_nik = ?", email, nip).First(&existing).Error == nil {
			skipped++
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			skipped++
			continue
		}

		user := models.User{
			FullName:     fullName,
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "teacher",
			NisNipNik:    nip,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			skipped++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import selesai",
		"imported": imported,
		"skipped":  skipped,
	})
}

func (h *ImportHandler) ImportDudi(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can import"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV empty or invalid"})
		return
	}

	colIndex := mapCSVHeaders(headers)
	required := []string{"company_name", "address"}
	for _, r := range required {
		if _, ok := colIndex[r]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom '" + r + "' wajib ada di CSV"})
			return
		}
	}

	imported := 0
	skipped := 0

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		companyName := getCol(record, colIndex, "company_name")
		address := getCol(record, colIndex, "address")
		picName := getCol(record, colIndex, "pic_name")
		phone := getCol(record, colIndex, "phone")
		radiusStr := getCol(record, colIndex, "radius_allowed")
		latStr := getCol(record, colIndex, "latitude")
		lngStr := getCol(record, colIndex, "longitude")

		if companyName == "" {
			skipped++
			continue
		}

		radius := 500
		if radiusStr != "" {
			if v, err := strconv.Atoi(radiusStr); err == nil {
				radius = v
			}
		}

		lat := 0.0
		if latStr != "" {
			if v, err := strconv.ParseFloat(latStr, 64); err == nil {
				lat = v
			}
		}

		lng := 0.0
		if lngStr != "" {
			if v, err := strconv.ParseFloat(lngStr, 64); err == nil {
				lng = v
			}
		}

		dudi := models.DUDI{
			CompanyName:   companyName,
			Address:       address,
			Latitude:      lat,
			Longitude:     lng,
			RadiusAllowed: radius,
			PicName:       picName,
			Phone:         phone,
		}

		dudi.ID = uuid.New()
		if err := database.DB.Create(&dudi).Error; err != nil {
			skipped++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import selesai",
		"imported": imported,
		"skipped":  skipped,
	})
}

func (h *ImportHandler) ImportInstrukturDudi(c *gin.Context) {
	role, _ := c.Get("role")
	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can import"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File CSV required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CSV empty or invalid"})
		return
	}

	colIndex := mapCSVHeaders(headers)
	required := []string{"full_name", "email", "nik", "dudi_nik"}
	for _, r := range required {
		if _, ok := colIndex[r]; !ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kolom '" + r + "' wajib ada di CSV"})
			return
		}
	}

	imported := 0
	skipped := 0
	errors := []string{}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}

		fullName := getCol(record, colIndex, "full_name")
		email := getCol(record, colIndex, "email")
		nik := getCol(record, colIndex, "nik")
		password := getCol(record, colIndex, "password")
		dudiNIK := getCol(record, colIndex, "dudi_nik")

		if fullName == "" || email == "" || nik == "" || dudiNIK == "" {
			skipped++
			continue
		}

		if password == "" {
			password = "pkl123456"
		}

		var existing models.User
		if database.DB.Where("email = ? OR nis_nip_nik = ?", email, nik).First(&existing).Error == nil {
			skipped++
			continue
		}

		var dudiUser models.User
		if database.DB.Where("nis_nip_nik = ? AND role = 'dudi'", dudiNIK).First(&dudiUser).Error != nil {
			errors = append(errors, "DUDI dengan NIK '"+dudiNIK+"' tidak ditemukan untuk user "+fullName)
			skipped++
			continue
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			skipped++
			continue
		}

		user := models.User{
			FullName:     fullName,
			Email:        email,
			PasswordHash: string(hashedPassword),
			Role:         "dudi",
			NisNipNik:    nik,
			DudiID:       dudiUser.DudiID,
		}

		if err := database.DB.Create(&user).Error; err != nil {
			skipped++
			continue
		}
		imported++
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Import selesai",
		"imported": imported,
		"skipped":  skipped,
		"errors":   errors,
	})
}

func mapCSVHeaders(headers []string) map[string]int {
	m := make(map[string]int)
	for i, h := range headers {
		m[strings.TrimSpace(strings.ToLower(h))] = i
	}
	return m
}

func getCol(record []string, colIndex map[string]int, col string) string {
	if idx, ok := colIndex[col]; ok && idx < len(record) {
		return strings.TrimSpace(record[idx])
	}
	return ""
}
