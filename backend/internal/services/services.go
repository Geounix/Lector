package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo interface {
		DB() *sql.DB
		Redis() *redis.Client
		CreateScanJob(ctx context.Context, libraryID, userID int64) (string, error)
		GetScanJob(ctx context.Context, scanID string) (status string, chaptersFound int, errorMsg *string, startedAt, completedAt *string, err error)
		UpdateScanJob(ctx context.Context, scanID, status string, chaptersFound int, errorMsg *string) error
		UpdateScanJobStarted(ctx context.Context, scanID string) error
		PublishScanCommand(ctx context.Context, libraryID int64, scanID string) error
	}
}

func NewService(repo interface {
	DB() *sql.DB
	Redis() *redis.Client
	CreateScanJob(ctx context.Context, libraryID, userID int64) (string, error)
	GetScanJob(ctx context.Context, scanID string) (status string, chaptersFound int, errorMsg *string, startedAt, completedAt *string, err error)
	UpdateScanJob(ctx context.Context, scanID, status string, chaptersFound int, errorMsg *string) error
	UpdateScanJobStarted(ctx context.Context, scanID string) error
	PublishScanCommand(ctx context.Context, libraryID int64, scanID string) error
}) *Service {
	return &Service{repo: repo}
}

type User struct {
	ID               int64
	Username         string
	Email            string
	PasswordHash     string
	Avatar           *string
	Role             string
	Theme            string
	Language         string
	TwoFactorEnabled bool
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, avatar, role, theme, language, 
			   two_factor_enabled, created_at, updated_at
		FROM users WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Avatar, &user.Role, &user.Theme, &user.Language,
		&user.TwoFactorEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	var user User
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, username, email, password_hash, avatar, role, theme, language, 
			   two_factor_enabled, created_at, updated_at
		FROM users WHERE id = $1
	`, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.Avatar, &user.Role, &user.Theme, &user.Language,
		&user.TwoFactorEnabled, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (s *Service) CreateUser(ctx context.Context, username, email, password string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var user User
	err = s.repo.DB().QueryRowContext(ctx, `
		INSERT INTO users (username, email, password_hash, role, theme, language)
		VALUES ($1, $2, $3, 'user', 'dark', 'en')
		RETURNING id, username, email, role, theme, language, created_at, updated_at
	`, username, email, string(hashedPassword)).Scan(
		&user.ID, &user.Username, &user.Email, &user.Role, &user.Theme, &user.Language,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) ValidatePassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}

func (s *Service) GetLibraries(ctx context.Context) ([]Library, error) {
	rows, err := s.repo.DB().QueryContext(ctx, `
		SELECT id, name, type, path, watch_enabled, scan_interval, created_at
		FROM library ORDER BY name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var libraries []Library
	for rows.Next() {
		var lib Library
		if err := rows.Scan(&lib.ID, &lib.Name, &lib.Type, &lib.Path, &lib.WatchEnabled, &lib.ScanInterval, &lib.CreatedAt); err != nil {
			return nil, err
		}
		libraries = append(libraries, lib)
	}
	return libraries, nil
}

func (s *Service) CreateLibrary(ctx context.Context, name, libType, path string) (*Library, error) {
	var lib Library
	err := s.repo.DB().QueryRowContext(ctx, `
		INSERT INTO library (name, type, path, watch_enabled, scan_interval)
		VALUES ($1, $2, $3, true, 3600)
		RETURNING id, name, type, path, watch_enabled, scan_interval, created_at
	`, name, libType, path).Scan(
		&lib.ID, &lib.Name, &lib.Type, &lib.Path, &lib.WatchEnabled, &lib.ScanInterval, &lib.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &lib, nil
}

func (s *Service) GetSeries(ctx context.Context, libraryID *int64) ([]Series, error) {
	query := `
		SELECT id, library_id, title, sort_title, description, cover, year, status, 
			   publisher, language, age_rating, metadata_source, created_at, updated_at
		FROM series`
	var args []interface{}
	if libraryID != nil {
		query += " WHERE library_id = $1"
		args = append(args, *libraryID)
	}
	query += " ORDER BY title"

	rows, err := s.repo.DB().QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.LibraryID, &s.Title, &s.SortTitle, &s.Description,
			&s.Cover, &s.Year, &s.Status, &s.Publisher, &s.Language, &s.AgeRating,
			&s.MetadataSource, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		seriesList = append(seriesList, s)
	}
	return seriesList, nil
}

func (s *Service) GetChaptersBySeries(ctx context.Context, seriesID int64) ([]Chapter, error) {
	rows, err := s.repo.DB().QueryContext(ctx, `
		SELECT c.id, c.volume_id, c.title, c.file_path, c.hash, c.page_count, c.size, c.created_at
		FROM chapter c
		JOIN volume v ON c.volume_id = v.id
		WHERE v.series_id = $1
		ORDER BY v.number, c.title
	`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chapters []Chapter
	for rows.Next() {
		var ch Chapter
		if err := rows.Scan(&ch.ID, &ch.VolumeID, &ch.Title, &ch.FilePath, &ch.Hash, &ch.PageCount, &ch.Size, &ch.CreatedAt); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	return chapters, nil
}

func (s *Service) GetChapter(ctx context.Context, chapterID int64) (*Chapter, error) {
	var ch Chapter
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, volume_id, title, file_path, hash, page_count, size, created_at
		FROM chapter WHERE id = $1
	`, chapterID).Scan(&ch.ID, &ch.VolumeID, &ch.Title, &ch.FilePath, &ch.Hash, &ch.PageCount, &ch.Size, &ch.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &ch, nil
}

func (s *Service) GetReadingProgress(ctx context.Context, userID, chapterID int64) (*ReadingProgress, error) {
	var rp ReadingProgress
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, user_id, chapter_id, page, percentage, time_spent, last_read
		FROM reading_progress
		WHERE user_id = $1 AND chapter_id = $2
	`, userID, chapterID).Scan(&rp.ID, &rp.UserID, &rp.ChapterID, &rp.Page, &rp.Percentage, &rp.TimeSpent, &rp.LastRead)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rp, nil
}

func (s *Service) UpsertReadingProgress(ctx context.Context, userID, chapterID int64, page int, percentage float64) error {
	_, err := s.repo.DB().ExecContext(ctx, `
		INSERT INTO reading_progress (user_id, chapter_id, page, percentage, last_read)
		VALUES ($1, $2, $3, $4, NOW())
		ON CONFLICT (user_id, chapter_id) 
		DO UPDATE SET page = $3, percentage = $4, last_read = NOW()
	`, userID, chapterID, page, percentage)
	return err
}

func (s *Service) SearchSeries(ctx context.Context, query string) ([]Series, error) {
	rows, err := s.repo.DB().QueryContext(ctx, `
		SELECT id, library_id, title, sort_title, description, cover, year, status, 
			   publisher, language, age_rating, metadata_source, created_at, updated_at
		FROM series
		WHERE title ILIKE $1 OR sort_title ILIKE $1
		ORDER BY title
		LIMIT 50
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seriesList []Series
	for rows.Next() {
		var s Series
		if err := rows.Scan(&s.ID, &s.LibraryID, &s.Title, &s.SortTitle, &s.Description,
			&s.Cover, &s.Year, &s.Status, &s.Publisher, &s.Language, &s.AgeRating,
			&s.MetadataSource, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		seriesList = append(seriesList, s)
	}
	return seriesList, nil
}

func (s *Service) GetStats(ctx context.Context, userID int64) (*Stats, error) {
	var stats Stats

	s.repo.DB().QueryRowContext(ctx, `
		SELECT COUNT(*) FROM reading_progress WHERE user_id = $1
	`, userID).Scan(&stats.TotalChaptersRead)

	s.repo.DB().QueryRowContext(ctx, `
		SELECT COALESCE(SUM(time_spent), 0) FROM reading_progress WHERE user_id = $1
	`, userID).Scan(&stats.TotalReadingTime)

	s.repo.DB().QueryRowContext(ctx, `
		SELECT COUNT(DISTINCT v.series_id)
		FROM reading_progress rp
		JOIN chapter c ON rp.chapter_id = c.id
		JOIN volume v ON c.volume_id = v.id
		WHERE rp.user_id = $1 AND rp.percentage >= 100
	`, userID).Scan(&stats.TotalSeriesCompleted)

	return &stats, nil
}

func (s *Service) UpdateLibrary(ctx context.Context, id int64, name, libType, path string, watchEnabled bool, scanInterval int) (*Library, error) {
	var lib Library
	err := s.repo.DB().QueryRowContext(ctx, `
		UPDATE library SET name = $2, type = $3, path = $4, watch_enabled = $5, scan_interval = $6
		WHERE id = $1
		RETURNING id, name, type, path, watch_enabled, scan_interval, created_at
	`, id, name, libType, path, watchEnabled, scanInterval).Scan(
		&lib.ID, &lib.Name, &lib.Type, &lib.Path, &lib.WatchEnabled, &lib.ScanInterval, &lib.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &lib, nil
}

func (s *Service) DeleteLibrary(ctx context.Context, id int64) error {
	_, err := s.repo.DB().ExecContext(ctx, `DELETE FROM library WHERE id = $1`, id)
	return err
}

func (s *Service) GetSeriesByID(ctx context.Context, id int64) (*Series, []Volume, []Chapter, error) {
	var series Series
	err := s.repo.DB().QueryRowContext(ctx, `
		SELECT id, library_id, title, sort_title, description, cover, year, status,
			   publisher, language, age_rating, metadata_source, created_at, updated_at
		FROM series WHERE id = $1
	`, id).Scan(
		&series.ID, &series.LibraryID, &series.Title, &series.SortTitle, &series.Description,
		&series.Cover, &series.Year, &series.Status, &series.Publisher, &series.Language,
		&series.AgeRating, &series.MetadataSource, &series.CreatedAt, &series.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}

	volumes, err := s.GetVolumesBySeries(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	chapters, err := s.GetChaptersBySeries(ctx, id)
	if err != nil {
		return nil, nil, nil, err
	}

	return &series, volumes, chapters, nil
}

func (s *Service) GetVolumesBySeries(ctx context.Context, seriesID int64) ([]Volume, error) {
	rows, err := s.repo.DB().QueryContext(ctx, `
		SELECT id, series_id, number, title FROM volume WHERE series_id = $1 ORDER BY number
	`, seriesID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var volumes []Volume
	for rows.Next() {
		var v Volume
		if err := rows.Scan(&v.ID, &v.SeriesID, &v.Number, &v.Title); err != nil {
			return nil, err
		}
		volumes = append(volumes, v)
	}
	return volumes, nil
}

func (s *Service) TriggerScan(ctx context.Context, libraryID, userID int64) (string, error) {
	scanID, err := s.repo.CreateScanJob(ctx, libraryID, userID)
	if err != nil {
		return "", err
	}

	err = s.repo.PublishScanCommand(ctx, libraryID, scanID)
	if err != nil {
		return "", err
	}

	return scanID, nil
}

func (s *Service) GetScanStatus(ctx context.Context, scanID string) (*ScanStatus, error) {
	status, chaptersFound, errorMsg, startedAt, completedAt, err := s.repo.GetScanJob(ctx, scanID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	scanStatus := &ScanStatus{
		ScanID:         scanID,
		Status:         status,
		ChaptersFound:  chaptersFound,
		ErrorMessage:   errorMsg,
	}

	if startedAt != nil {
		scanStatus.StartedAt = *startedAt
	}
	if completedAt != nil {
		scanStatus.CompletedAt = *completedAt
	}

	return scanStatus, nil
}

type Volume struct {
	ID       int64   `json:"id"`
	SeriesID int64   `json:"series_id"`
	Number   int     `json:"number"`
	Title    *string `json:"title"`
}

type ScanStatus struct {
	ScanID         string  `json:"scan_id"`
	Status         string  `json:"status"`
	ChaptersFound  int     `json:"chapters_found"`
	ErrorMessage  *string `json:"error_message,omitempty"`
	StartedAt      string  `json:"started_at,omitempty"`
	CompletedAt    string  `json:"completed_at,omitempty"`
}

type Library struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	Path         string    `json:"path"`
	WatchEnabled bool      `json:"watch_enabled"`
	ScanInterval int       `json:"scan_interval"`
	CreatedAt    time.Time `json:"created_at"`
}

type Series struct {
	ID             int64      `json:"id"`
	LibraryID      int64      `json:"library_id"`
	Title          string     `json:"title"`
	SortTitle      *string    `json:"sort_title"`
	Description    *string     `json:"description"`
	Cover          *string    `json:"cover"`
	Year           *int       `json:"year"`
	Status         *string    `json:"status"`
	Publisher      *string    `json:"publisher"`
	Language       *string    `json:"language"`
	AgeRating      *string    `json:"age_rating"`
	MetadataSource *string    `json:"metadata_source"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type Chapter struct {
	ID        int64     `json:"id"`
	VolumeID  int64     `json:"volume_id"`
	Title     *string   `json:"title"`
	FilePath  string    `json:"file_path"`
	Hash      string    `json:"hash"`
	PageCount int       `json:"page_count"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type ReadingProgress struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ChapterID  int64     `json:"chapter_id"`
	Page       int       `json:"page"`
	Percentage float64   `json:"percentage"`
	TimeSpent  int64     `json:"time_spent"`
	LastRead   time.Time `json:"last_read"`
}

type Stats struct {
	TotalChaptersRead   int   `json:"total_chapters_read"`
	TotalReadingTime     int64 `json:"total_reading_time"`
	TotalSeriesCompleted int   `json:"total_series_completed"`
}