CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    avatar TEXT,
    role VARCHAR(20) DEFAULT 'user',
    theme VARCHAR(20) DEFAULT 'dark',
    language VARCHAR(10) DEFAULT 'en',
    two_factor_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS library (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(20) DEFAULT 'comics',
    path TEXT NOT NULL,
    watch_enabled BOOLEAN DEFAULT TRUE,
    scan_interval INTEGER DEFAULT 3600,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS series (
    id BIGSERIAL PRIMARY KEY,
    library_id BIGINT REFERENCES library(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    sort_title TEXT,
    description TEXT,
    cover TEXT,
    year INTEGER,
    status VARCHAR(20),
    publisher TEXT,
    language VARCHAR(10),
    age_rating VARCHAR(10),
    metadata_source VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS volume (
    id BIGSERIAL PRIMARY KEY,
    series_id BIGINT REFERENCES series(id) ON DELETE CASCADE,
    number INTEGER NOT NULL,
    title TEXT
);

CREATE TABLE IF NOT EXISTS chapter (
    id BIGSERIAL PRIMARY KEY,
    volume_id BIGINT REFERENCES volume(id) ON DELETE CASCADE,
    title TEXT,
    file_path TEXT NOT NULL,
    hash VARCHAR(64) UNIQUE NOT NULL,
    page_count INTEGER DEFAULT 0,
    size BIGINT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reading_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    chapter_id BIGINT REFERENCES chapter(id) ON DELETE CASCADE,
    page INTEGER DEFAULT 0,
    percentage NUMERIC(5,2) DEFAULT 0,
    time_spent BIGINT DEFAULT 0,
    last_read TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, chapter_id)
);

CREATE TABLE IF NOT EXISTS collection (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
    is_smart BOOLEAN DEFAULT FALSE,
    query JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS collection_item (
    id BIGSERIAL PRIMARY KEY,
    collection_id BIGINT REFERENCES collection(id) ON DELETE CASCADE,
    chapter_id BIGINT REFERENCES chapter(id) ON DELETE CASCADE,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(collection_id, chapter_id)
);

CREATE TABLE IF NOT EXISTS metadata (
    id BIGSERIAL PRIMARY KEY,
    series_id BIGINT UNIQUE REFERENCES series(id) ON DELETE CASCADE,
    writer TEXT,
    artist TEXT,
    genres TEXT[],
    tags TEXT[],
    isbn TEXT,
    rating NUMERIC(3,1),
    release_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_reading_progress_user ON reading_progress(user_id);
CREATE INDEX idx_reading_progress_chapter ON reading_progress(chapter_id);
CREATE INDEX idx_scan_job_library ON scan_job(library_id);
CREATE INDEX idx_scan_job_user ON scan_job(user_id);
CREATE INDEX idx_scan_job_status ON scan_job(status);

CREATE TABLE IF NOT EXISTS scan_job (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    library_id BIGINT REFERENCES library(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(20) DEFAULT 'queued',
    chapters_found INT DEFAULT 0,
    error_message TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);