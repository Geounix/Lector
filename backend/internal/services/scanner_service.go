package services

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ScannerService struct {
	repo *Repository
}

func NewScannerService(repo *Repository) *ScannerService {
	return &ScannerService{repo: repo}
}

type LibraryInfo struct {
	ID           int64
	Name         string
	Path         string
	WatchEnabled bool
}

func (s *ScannerService) ScanLibrary(ctx context.Context, libraryID int64) error {
	var lib LibraryInfo
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, name, path, watch_enabled FROM library WHERE id = $1
	`, libraryID).Scan(&lib.ID, &lib.Name, &lib.Path, &lib.WatchEnabled)
	if err != nil {
		return fmt.Errorf("library not found: %w", err)
	}

	return s.scanDirectory(ctx, &lib)
}

func (s *ScannerService) scanDirectory(ctx context.Context, lib *LibraryInfo) error {
	info, err := os.Stat(lib.Path)
	if err != nil {
		return fmt.Errorf("path error: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("path is not a directory")
	}

	entries, err := os.ReadDir(lib.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		seriesPath := filepath.Join(lib.Path, entry.Name())
		seriesID, err := s.getOrCreateSeries(ctx, lib.ID, entry.Name())
		if err != nil {
			continue
		}

		w := &seriesWalker{
			repo:     s.repo,
			seriesID: seriesID,
		}
		w.walk(seriesPath)
	}

	return nil
}

type seriesWalker struct {
	repo     *Repository
	seriesID int64
}

func (w *seriesWalker) walk(dir string) {
	entries, _ := os.ReadDir(dir)
	for _, entry := range entries {
		if entry.IsDir() {
			w.walk(filepath.Join(dir, entry.Name()))
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext != ".cbz" && ext != ".zip" && ext != ".cbr" {
			continue
		}

		w.processFile(filepath.Join(dir, entry.Name()))
	}
}

func (w *seriesWalker) processFile(filePath string) {
	filename := filepath.Base(filePath)
	hash := hashFile(filePath)
	pageCount := countZipPages(filePath)
	volumeNum := extractVolumeNumber(filename)

	volumeID, _ := w.getOrCreateVolume(volumeNum)

	_, _ = w.repo.DB().ExecContext(context.Background(), `
		INSERT INTO chapter (volume_id, title, file_path, hash, page_count, size)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (hash) DO NOTHING
	`, volumeID, filename, filePath, hash, pageCount, getFileSize(filePath))
}

func (w *seriesWalker) getOrCreateVolume(number int) (int64, error) {
	var existingID int64
	err := w.repo.DB().QueryRowContext(context.Background(), `
		SELECT id FROM volume WHERE series_id = $1 AND number = $2
	`, w.seriesID, number).Scan(&existingID)
	if err == nil {
		return existingID, nil
	}

	var newID int64
	err = w.repo.DB().QueryRowContext(context.Background(), `
		INSERT INTO volume (series_id, number) VALUES ($1, $2) RETURNING id
	`, w.seriesID, number).Scan(&newID)
	return newID, err
}

func (s *ScannerService) getOrCreateSeries(ctx context.Context, libraryID int64, title string) (int64, error) {
	var existingID int64
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id FROM series WHERE library_id = $1 AND title = $2
	`, libraryID, title).Scan(&existingID)
	if err == nil {
		return existingID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	var newID int64
	err = s.repo.DB().QueryRowContext(ctx, `
		INSERT INTO series (library_id, title, sort_title)
		VALUES ($1, $2, $3)
		RETURNING id
	`, libraryID, title, generateSortTitle(title)).Scan(&newID)
	return newID, err
}

func generateSortTitle(title string) string {
	sortTitle := strings.ToLower(title)
	sortTitle = strings.ReplaceAll(sortTitle, "the ", "")
	sortTitle = strings.ReplaceAll(sortTitle, "a ", "")
	sortTitle = strings.ReplaceAll(sortTitle, "an ", "")
	return strings.TrimSpace(sortTitle)
}

func hashFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil || len(data) < 1024 {
		return ""
	}

	h := sha256.New()
	h.Write(data[:1024])
	return fmt.Sprintf("%x", h.Sum(nil))
}

func countZipPages(filePath string) int {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext != ".cbz" && ext != ".zip" {
		return 0
	}

	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return 0
	}
	defer reader.Close()

	count := 0
	for _, file := range reader.File {
		if isImageFile(file.Name) {
			count++
		}
	}
	return count
}

func extractVolumeNumber(filename string) int {
	re := regexp.MustCompile(`(?i)vol\.?\s*(\d+)`)
	matches := re.FindStringSubmatch(filename)
	if len(matches) > 1 {
		if v, err := strconv.Atoi(matches[1]); err == nil {
			return v
		}
	}
	return 1
}

func getFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif"
}

func (s *ScannerService) ScanAllLibraries(ctx context.Context) error {
	rows, err := s.repo.DB().QueryContext(ctx, `
		SELECT id, name, path, watch_enabled FROM library WHERE watch_enabled = true
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var lib LibraryInfo
		if err := rows.Scan(&lib.ID, &lib.Name, &lib.Path, &lib.WatchEnabled); err != nil {
			continue
		}
		go s.ScanLibrary(ctx, lib.ID)
	}

	return nil
}

var _ = time.Time{}
var _ = uuid.New