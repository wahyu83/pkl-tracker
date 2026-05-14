package storage

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type DriveStorage struct {
	service   *drive.Service
	folderID  string
	enabled   bool
	localPath string
}

func NewDriveStorage(credentialsFile, folderID string) *DriveStorage {
	s := &DriveStorage{
		folderID:  folderID,
		localPath: "uploads",
	}

	if credentialsFile == "" {
		log.Println("[Storage] Google Drive disabled, using local storage at:", s.localPath)
		return s
	}

	if _, err := os.Stat(credentialsFile); os.IsNotExist(err) {
		log.Println("[Storage] Service account file not found, using local storage at:", s.localPath)
		return s
	}

	ctx := context.Background()
	service, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		log.Printf("[Storage] Failed to init Google Drive: %v, using local storage", err)
		return s
	}

	s.service = service
	s.enabled = true
	log.Println("[Storage] Google Drive connected, folder:", folderID)
	return s
}

func (s *DriveStorage) UploadFile(file io.Reader, filename, userID string) (string, error) {
	if s.enabled {
		return s.uploadToDrive(file, filename, userID)
	}
	return s.saveLocal(file, filename, userID)
}

func (s *DriveStorage) uploadToDrive(file io.Reader, filename, userID string) (string, error) {
	month := time.Now().Format("2006-01")
	folderPath := fmt.Sprintf("pkl_photos/%s/%s", userID, month)

	parentID := s.findOrCreateFolder(folderPath)

	driveFile := &drive.File{
		Name:     filename,
		Parents:  []string{parentID},
		MimeType: mime.TypeByExtension(filepath.Ext(filename)),
	}

	created, err := s.service.Files.Create(driveFile).Media(file).Do()
	if err != nil {
		return "", fmt.Errorf("upload to drive failed: %w", err)
	}

	_, err = s.service.Permissions.Create(created.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()
	if err != nil {
		log.Printf("[Storage] Failed to set public permission: %v", err)
	}

	_, err = s.service.Files.Update(created.Id, &drive.File{
		Description: "Uploaded via PKL Tracker",
	}).Do()
	if err != nil {
		log.Printf("[Storage] Failed to update file description: %v", err)
	}

	return fmt.Sprintf("https://drive.google.com/file/d/%s/view", created.Id), nil
}

func (s *DriveStorage) findOrCreateFolder(folderPath string) string {
	currentParent := s.folderID

	parts := []string{}
	remaining := folderPath
	for remaining != "" {
		dir := filepath.Dir(remaining)
		if dir == "." || dir == remaining {
			parts = append([]string{remaining}, parts...)
			break
		}
		name := filepath.Base(remaining)
		parts = append([]string{name}, parts...)
		remaining = dir
	}

	for _, folderName := range parts {
		query := fmt.Sprintf("'%s' in parents and name = '%s' and mimeType = 'application/vnd.google-apps.folder' and trashed = false",
			currentParent, folderName)

		files, err := s.service.Files.List().Q(query).PageSize(1).Do()
		if err == nil && len(files.Files) > 0 {
			currentParent = files.Files[0].Id
			continue
		}

		newFolder := &drive.File{
			Name:     folderName,
			MimeType: "application/vnd.google-apps.folder",
			Parents:  []string{currentParent},
		}

		created, err := s.service.Files.Create(newFolder).Do()
		if err != nil {
			return currentParent
		}
		currentParent = created.Id
	}

	return currentParent
}

func (s *DriveStorage) saveLocal(file io.Reader, filename, userID string) (string, error) {
	month := time.Now().Format("2006-01")
	dir := filepath.Join(s.localPath, userID, month)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	fullPath := filepath.Join(dir, filename)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return "/uploads/" + userID + "/" + month + "/" + filename, nil
}
