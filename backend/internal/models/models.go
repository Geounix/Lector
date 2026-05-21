package models

import (
	"time"
)

type User struct {
	ID              int64     `json:"id"`
	Username        string    `json:"username"`
	Email           string    `json:"email"`
	PasswordHash    string    `json:"-"`
	Avatar          *string   `json:"avatar"`
	Role            string    `json:"role"`
	Theme           string    `json:"theme"`
	Language        string    `json:"language"`
	TwoFactorEnabled bool     `json:"two_factor_enabled"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type Library struct {
	ID            int64     `json:"id"`
	Name          string    `json:"name"`
	Type          string    `json:"type"`
	Path          string    `json:"path"`
	WatchEnabled  bool      `json:"watch_enabled"`
	ScanInterval  int       `json:"scan_interval"`
	CreatedAt     time.Time `json:"created_at"`
}

type Series struct {
	ID            int64     `json:"id"`
	LibraryID     int64     `json:"library_id"`
	Title         string    `json:"title"`
	SortTitle     *string   `json:"sort_title"`
	Description   *string   `json:"description"`
	Cover         *string   `json:"cover"`
	Year          *int      `json:"year"`
	Status        *string   `json:"status"`
	Publisher     *string   `json:"publisher"`
	Language      *string   `json:"language"`
	AgeRating     *string   `json:"age_rating"`
	MetadataSource *string  `json:"metadata_source"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Volume struct {
	ID       int64  `json:"id"`
	SeriesID int64  `json:"series_id"`
	Number   int    `json:"number"`
	Title    string `json:"title"`
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

type Metadata struct {
	ID          int64      `json:"id"`
	SeriesID    int64      `json:"series_id"`
	Writer      *string    `json:"writer"`
	Artist      *string    `json:"artist"`
	Genres      []string   `json:"genres"`
	Tags        []string   `json:"tags"`
	ISBN        *string    `json:"isbn"`
	Rating      *float64   `json:"rating"`
	ReleaseDate *time.Time `json:"release_date"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type Collection struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	UserID    int64     `json:"user_id"`
	IsSmart   bool      `json:"is_smart"`
	Query     *string   `json:"query"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}