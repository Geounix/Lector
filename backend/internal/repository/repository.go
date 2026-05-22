package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/redis/go-redis/v9"
	_ "github.com/lib/pq"
)

type Repository struct {
	db    *sql.DB
	redis *redis.Client
}

func NewPostgresDB(connStr string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	return db, nil
}

func NewRedisClient(addr string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

func NewRepository(db *sql.DB, redis *redis.Client) *Repository {
	return &Repository{
		db:    db,
		redis: redis,
	}
}

func (r *Repository) DB() *sql.DB {
	return r.db
}

func (r *Repository) Redis() *redis.Client {
	return r.redis
}

func (r *Repository) CacheSet(ctx context.Context, key string, value interface{}) error {
	return r.redis.Set(ctx, key, value, 0).Err()
}

func (r *Repository) CacheGet(ctx context.Context, key string) (string, error) {
	return r.redis.Get(ctx, key).Result()
}

func (r *Repository) CacheDelete(ctx context.Context, key string) error {
	return r.redis.Del(ctx, key).Err()
}

func (r *Repository) Close() error {
	if err := r.db.Close(); err != nil {
		return fmt.Errorf("failed to close db: %w", err)
	}
	if err := r.redis.Close(); err != nil {
		return fmt.Errorf("failed to close redis: %w", err)
	}
	return nil
}

func (r *Repository) CreateScanJob(ctx context.Context, libraryID, userID int64) (string, error) {
	var id string
	err := r.db.QueryRowContext(ctx, `
		INSERT INTO scan_job (library_id, user_id, status)
		VALUES ($1, $2, 'queued')
		RETURNING id
	`, libraryID, userID).Scan(&id)
	return id, err
}

func (r *Repository) GetScanJob(ctx context.Context, scanID string) (status string, chaptersFound int, errorMsg *string, startedAt, completedAt *string, err error) {
	err = r.db.QueryRowContext(ctx, `
		SELECT status, chapters_found, error_message, started_at, completed_at
		FROM scan_job WHERE id = $1
	`, scanID).Scan(&status, &chaptersFound, &errorMsg, &startedAt, &completedAt)
	return
}

func (r *Repository) UpdateScanJob(ctx context.Context, scanID, status string, chaptersFound int, errorMsg *string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE scan_job
		SET status = $2, chapters_found = $3, error_message = $4, completed_at = CASE WHEN $2 IN ('completed', 'failed') THEN NOW() ELSE completed_at END
		WHERE id = $1
	`, scanID, status, chaptersFound, errorMsg)
	return err
}

func (r *Repository) UpdateScanJobStarted(ctx context.Context, scanID string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE scan_job SET status = 'scanning', started_at = NOW() WHERE id = $1
	`, scanID)
	return err
}

func (r *Repository) PublishScanCommand(ctx context.Context, libraryID int64, scanID string) error {
	msg := fmt.Sprintf(`{"library_id":%d,"scan_id":"%s"}`, libraryID, scanID)
	return r.redis.Publish(ctx, "lector:scan_library", msg).Err()
}