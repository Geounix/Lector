package workers

import (
	"github.com/lector-comics/lector/internal/repository"
)

type ScannerWorker struct {
	repo *repository.Repository
}

func NewScannerWorker(repo *repository.Repository) *ScannerWorker {
	return &ScannerWorker{repo: repo}
}

func (w *ScannerWorker) Start() {
}

type ThumbnailWorker struct {
	repo *repository.Repository
}

func NewThumbnailWorker(repo *repository.Repository) *ThumbnailWorker {
	return &ThumbnailWorker{repo: repo}
}

func (w *ThumbnailWorker) Start() {
}

type CleanupWorker struct {
	repo *repository.Repository
}

func NewCleanupWorker(repo *repository.Repository) *CleanupWorker {
	return &CleanupWorker{repo: repo}
}

func (w *CleanupWorker) Start() {
}