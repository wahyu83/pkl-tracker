package handlers

import (
	"compress/gzip"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"pkl-tracker/config"
)

type BackupHandler struct {
	cfg *config.Config
}

func NewBackupHandler(cfg *config.Config) *BackupHandler {
	return &BackupHandler{cfg: cfg}
}

func (h *BackupHandler) backupDir() string {
	if h.cfg.BackupDir == "" {
		return "backups"
	}
	return h.cfg.BackupDir
}

type BackupInfo struct {
	Filename    string    `json:"filename"`
	Size        int64     `json:"size"`
	SizeHuman   string    `json:"size_human"`
	CreatedAt   time.Time `json:"created_at"`
	DownloadURL string    `json:"download_url"`
}

func (h *BackupHandler) ListBackups(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	dir := h.backupDir()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		c.JSON(http.StatusOK, gin.H{"data": []BackupInfo{}})
		return
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read backup directory"})
		return
	}

	var backups []BackupInfo
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql.gz") {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		backups = append(backups, BackupInfo{
			Filename:    e.Name(),
			Size:        info.Size(),
			SizeHuman:   formatSize(info.Size()),
			CreatedAt:   info.ModTime(),
			DownloadURL: "/api/admin/backups/" + e.Name(),
		})
	}

	sort.Slice(backups, func(i, j int) bool {
		return backups[i].CreatedAt.After(backups[j].CreatedAt)
	})

	if backups == nil {
		backups = []BackupInfo{}
	}

	c.JSON(http.StatusOK, gin.H{"data": backups})
}

func (h *BackupHandler) CreateBackup(c *gin.Context) {
	isAdmin, _ := checkAdminAccess(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can create backup"})
		return
	}

	dir := h.backupDir()
	if err := os.MkdirAll(dir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup directory"})
		return
	}

	filename := fmt.Sprintf("%s_%s.sql.gz", h.cfg.DBName, time.Now().Format("20060102_150405"))
	dest := filepath.Join(dir, filename)

	cmd := exec.Command("pg_dump",
		"-h", h.cfg.DBHost,
		"-p", h.cfg.DBPort,
		"-U", h.cfg.DBUser,
		"-d", h.cfg.DBName,
		"--no-owner",
		"--no-privileges",
	)
	cmd.Env = append(os.Environ(), "PGPASSWORD="+h.cfg.DBPass)

	out, err := os.Create(dest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup file"})
		return
	}
	defer out.Close()

	gz := gzip.NewWriter(out)
	cmd.Stdout = gz
	var stderr strings.Builder
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		gz.Close()
		out.Close()
		os.Remove(dest)
		log.Printf("[Backup] pg_dump failed: %v | %s", err, stderr.String())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Backup failed: " + stderr.String()})
		return
	}

	if err := gz.Close(); err != nil {
		out.Close()
		os.Remove(dest)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Backup failed: gzip error"})
		return
	}

	info, err := os.Stat(dest)
	if err != nil || info.Size() == 0 {
		out.Close()
		os.Remove(dest)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Backup failed: empty output"})
		return
	}

	log.Printf("[Backup] Created %s (%s)", filename, formatSize(info.Size()))
	c.JSON(http.StatusCreated, gin.H{
		"message": "Backup created",
		"data": BackupInfo{
			Filename:    filename,
			Size:        info.Size(),
			SizeHuman:   formatSize(info.Size()),
			CreatedAt:   info.ModTime(),
			DownloadURL: "/api/admin/backups/" + filename,
		},
	})
}

func (h *BackupHandler) DownloadBackup(c *gin.Context) {
	isAdmin, adminJurusan := checkAdminAccess(c)
	if !isAdmin && adminJurusan == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	filename := filepath.Base(c.Param("filename"))
	if filename != c.Param("filename") || !strings.HasSuffix(filename, ".sql.gz") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	fullPath := filepath.Join(h.backupDir(), filename)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
		return
	}

	c.FileAttachment(fullPath, filename)
}

func (h *BackupHandler) DeleteBackup(c *gin.Context) {
	isAdmin, _ := checkAdminAccess(c)
	if !isAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only admin can delete backup"})
		return
	}

	filename := filepath.Base(c.Param("filename"))
	if filename != c.Param("filename") || !strings.HasSuffix(filename, ".sql.gz") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid filename"})
		return
	}

	fullPath := filepath.Join(h.backupDir(), filename)
	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Backup not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete backup"})
		return
	}

	log.Printf("[Backup] Deleted %s", filename)
	c.JSON(http.StatusOK, gin.H{"message": "Backup deleted"})
}

func formatSize(n int64) string {
	const unit = 1024
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := int64(unit), 0
	for m := n / unit; m >= unit; m /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(n)/float64(div), "KMGTPE"[exp])
}
