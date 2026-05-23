package workers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lector-comics/lector/internal/repository"
	"github.com/lector-comics/lector/internal/scanner"
)

type ScannerWorker struct {
	repo    *repository.Repository
	parser  *scanner.ComicParser
	scanCtx context.Context
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewScannerWorker(repo *repository.Repository) *ScannerWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &ScannerWorker{
		repo:   repo,
		parser: scanner.NewComicParser(),
		scanCtx: ctx,
		cancel:  cancel,
	}
}

func (w *ScannerWorker) Start() {
	log.Println("Scanner worker started")
	
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	w.scanAllLibraries()

	for {
		select {
		case <-w.scanCtx.Done():
			log.Println("Scanner worker stopping...")
			return
		case <-ticker.C:
			w.scanAllLibraries()
		}
	}
}

func (w *ScannerWorker) Stop() {
	w.cancel()
	w.wg.Wait()
}

func (w *ScannerWorker) SubscribeToScans() {
	log.Println("Scanner worker subscribing to scan commands")
	pubsub := w.repo.Redis().Subscribe(w.scanCtx, "lector:scan_library")
	defer pubsub.Close()

	ch := pubsub.Channel()
	log.Println("Scanner worker subscribed to channel 'lector:scan_library'")

	for {
		select {
		case <-w.scanCtx.Done():
			log.Println("Scanner worker unsubscribe from scans...")
			return
		case msg := <-ch:
			w.handleScanCommand(msg.Payload)
		}
	}
}

func (w *ScannerWorker) handleScanCommand(payload string) {
	var cmd ScanCommand
	if err := json.Unmarshal([]byte(payload), &cmd); err != nil {
		log.Printf("Scanner: failed to parse scan command: %v", err)
		return
	}

	log.Printf("Scanner: received scan command for library %d, scan_id: %s", cmd.LibraryID, cmd.ScanID)

	if err := w.repo.UpdateScanJobStarted(w.scanCtx, cmd.ScanID); err != nil {
		log.Printf("Scanner: failed to update scan job status: %v", err)
	}

	var chaptersFound int
	var errMsg *string

	scanFunc := func() error {
		rows, err := w.repo.DB().QueryContext(w.scanCtx, `
			SELECT id, name, path, watch_enabled FROM library WHERE id = $1
		`, cmd.LibraryID)
		if err != nil {
			return err
		}
		defer rows.Close()

		if !rows.Next() {
			return fmt.Errorf("library not found")
		}

		var lib LibraryInfo
		if err := rows.Scan(&lib.ID, &lib.Name, &lib.Path, &lib.WatchEnabled); err != nil {
			return err
		}

		return w.scanLibraryWithCount(w.scanCtx, &lib, &chaptersFound)
	}

	if err := scanFunc(); err != nil {
		errMsgStr := err.Error()
		errMsg = &errMsgStr
		w.repo.UpdateScanJob(w.scanCtx, cmd.ScanID, "failed", 0, errMsg)
		log.Printf("Scanner: scan %s failed: %v", cmd.ScanID, err)
		return
	}

	w.repo.UpdateScanJob(w.scanCtx, cmd.ScanID, "completed", chaptersFound, nil)
	log.Printf("Scanner: scan %s completed, found %d chapters", cmd.ScanID, chaptersFound)
}

func (w *ScannerWorker) scanLibraryWithCount(ctx context.Context, lib *LibraryInfo, chaptersFound *int) error {
	info, err := os.Stat(lib.Path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("library path does not exist: %s", lib.Path)
		}
		return fmt.Errorf("failed to stat library path: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("library path is not a directory: %s", lib.Path)
	}

	entries, err := os.ReadDir(lib.Path)
	if err != nil {
		return fmt.Errorf("failed to read library directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		seriesPath := filepath.Join(lib.Path, entry.Name())
		count, err := w.scanSeriesDirectoryWithCount(ctx, lib.ID, seriesPath, entry.Name())
		if err != nil {
			log.Printf("Scanner: failed to scan series directory: %v", err)
			continue
		}
		*chaptersFound += count
	}

	return nil
}

func (w *ScannerWorker) scanSeriesDirectoryWithCount(ctx context.Context, libraryID int64, seriesPath, seriesName string) (int, error) {
	entries, err := os.ReadDir(seriesPath)
	if err != nil {
		return 0, fmt.Errorf("failed to read series directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".cbz" || ext == ".cbr" || ext == ".zip" {
			files = append(files, filepath.Join(seriesPath, entry.Name()))
		}
	}

	if len(files) == 0 {
		return 0, nil
	}

	seriesID, err := w.getOrCreateSeries(ctx, libraryID, seriesName)
	if err != nil {
		return 0, fmt.Errorf("failed to get or create series: %w", err)
	}

	count := 0
	for _, filePath := range files {
		if w.processComicFileWithCount(ctx, seriesID, filePath) {
			count++
		}
	}

	return count, nil
}

func (w *ScannerWorker) processComicFileWithCount(ctx context.Context, seriesID int64, filePath string) bool {
	info := w.parser.ParseFilename(filePath)

	hash, err := scanner.CalculateHashPartial(filePath)
	if err != nil {
		log.Printf("Scanner: failed to calculate hash for %s: %v", filePath, err)
		return false
	}

	pageCount, _ := scanner.ExtractPageCount(filePath)

	volumeNumber := 1
	if info.Volume != nil {
		volumeNumber = *info.Volume
	}

	volumeID, err := w.getOrCreateVolume(ctx, seriesID, volumeNumber)
	if err != nil {
		log.Printf("Scanner: failed to get or create volume: %v", err)
		return false
	}

	result, err := w.repo.DB().ExecContext(ctx, `
		INSERT INTO chapter (volume_id, title, file_path, hash, page_count, size)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (hash) DO NOTHING
	`, volumeID, info.Title, filePath, hash, pageCount, scanner.GetFileSize(filePath))

	if err != nil {
		log.Printf("Scanner: failed to insert chapter: %v", err)
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

type ScanCommand struct {
	LibraryID int64  `json:"library_id"`
	ScanID    string `json:"scan_id"`
}

func (w *ScannerWorker) scanAllLibraries() {
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		
		ctx := context.Background()
		
		rows, err := w.repo.DB().QueryContext(ctx, `
			SELECT id, name, path, watch_enabled 
			FROM library 
			WHERE watch_enabled = true
		`)
		if err != nil {
			log.Printf("Scanner: failed to get libraries: %v", err)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var lib LibraryInfo
			if err := rows.Scan(&lib.ID, &lib.Name, &lib.Path, &lib.WatchEnabled); err != nil {
				log.Printf("Scanner: failed to scan row: %v", err)
				continue
			}

			log.Printf("Scanner: scanning library '%s' at %s", lib.Name, lib.Path)
			w.scanLibrary(ctx, &lib)
		}
	}()
}

func (w *ScannerWorker) scanLibrary(ctx context.Context, lib *LibraryInfo) {
	info, err := os.Stat(lib.Path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Scanner: library path does not exist: %s", lib.Path)
			return
		}
		log.Printf("Scanner: failed to stat library path: %v", err)
		return
	}

	if !info.IsDir() {
		log.Printf("Scanner: library path is not a directory: %s", lib.Path)
		return
	}

	entries, err := os.ReadDir(lib.Path)
	if err != nil {
		log.Printf("Scanner: failed to read library directory: %v", err)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		seriesPath := filepath.Join(lib.Path, entry.Name())
		w.scanSeriesDirectory(ctx, lib.ID, seriesPath, entry.Name())
	}
}

func (w *ScannerWorker) scanSeriesDirectory(ctx context.Context, libraryID int64, seriesPath, seriesName string) {
	entries, err := os.ReadDir(seriesPath)
	if err != nil {
		log.Printf("Scanner: failed to read series directory: %v", err)
		return
	}

	var files []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".cbz" || ext == ".cbr" || ext == ".zip" {
			files = append(files, filepath.Join(seriesPath, entry.Name()))
		}
	}

	if len(files) == 0 {
		return
	}

	seriesID, err := w.getOrCreateSeries(ctx, libraryID, seriesName)
	if err != nil {
		log.Printf("Scanner: failed to get or create series: %v", err)
		return
	}

	for _, filePath := range files {
		w.processComicFile(ctx, seriesID, filePath)
	}
}

func (w *ScannerWorker) processComicFile(ctx context.Context, seriesID int64, filePath string) {
	info := w.parser.ParseFilename(filePath)
	
	hash, err := scanner.CalculateHashPartial(filePath)
	if err != nil {
		log.Printf("Scanner: failed to calculate hash for %s: %v", filePath, err)
		return
	}

	pageCount, _ := scanner.ExtractPageCount(filePath)

	volumeNumber := 1
	if info.Volume != nil {
		volumeNumber = *info.Volume
	}

	volumeID, err := w.getOrCreateVolume(ctx, seriesID, volumeNumber)
	if err != nil {
		log.Printf("Scanner: failed to get or create volume: %v", err)
		return
	}

	_, err = w.repo.DB().ExecContext(ctx, `
		INSERT INTO chapter (volume_id, title, file_path, hash, page_count, size)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (hash) DO NOTHING
	`, volumeID, info.Title, filePath, hash, pageCount, scanner.GetFileSize(filePath))

	if err != nil {
		log.Printf("Scanner: failed to insert chapter: %v", err)
	}
}

func (w *ScannerWorker) getOrCreateSeries(ctx context.Context, libraryID int64, title string) (int64, error) {
	var existingID int64
	err := w.repo.DB().QueryRowContext(ctx, `
		SELECT id FROM series WHERE library_id = $1 AND title = $2
	`, libraryID, title).Scan(&existingID)

	if err == nil {
		return existingID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	var newID int64
	err = w.repo.DB().QueryRowContext(ctx, `
		INSERT INTO series (library_id, title, sort_title)
		VALUES ($1, $2, $3)
		RETURNING id
	`, libraryID, title, scanner.GenerateSortTitle(title)).Scan(&newID)

	return newID, err
}

func (w *ScannerWorker) getOrCreateVolume(ctx context.Context, seriesID int64, number int) (int64, error) {
	var existingID int64
	err := w.repo.DB().QueryRowContext(ctx, `
		SELECT id FROM volume WHERE series_id = $1 AND number = $2
	`, seriesID, number).Scan(&existingID)

	if err == nil {
		return existingID, nil
	}

	if err != sql.ErrNoRows {
		return 0, err
	}

	var newID int64
	err = w.repo.DB().QueryRowContext(ctx, `
		INSERT INTO volume (series_id, number)
		VALUES ($1, $2)
		RETURNING id
	`, seriesID, number).Scan(&newID)

	return newID, err
}

type LibraryInfo struct {
	ID           int64
	Name         string
	Path         string
	WatchEnabled bool
}

func EnsureDirectories(root string) error {
	dirs := []string{
		"library",
		"metadata",
		"cache",
		"backups",
		"logs",
	}

	for _, dir := range dirs {
		path := filepath.Join(root, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

var _ = uuid.New