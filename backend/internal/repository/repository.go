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