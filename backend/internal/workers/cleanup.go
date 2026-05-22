package workers

import (
	"context"
	"log"
	"time"

	"github.com/lector-comics/lector/internal/repository"
)

type CleanupWorker struct {
	repo    *repository.Repository
	scanCtx context.Context
	cancel  context.CancelFunc
}

func NewCleanupWorker(repo *repository.Repository) *CleanupWorker {
	ctx, cancel := context.WithCancel(context.Background())
	return &CleanupWorker{
		repo:    repo,
		scanCtx: ctx,
		cancel:  cancel,
	}
}

func (w *CleanupWorker) Start() {
	log.Println("Cleanup worker started")

	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-w.scanCtx.Done():
			log.Println("Cleanup worker stopping...")
			return
		case <-ticker.C:
			w.cleanup()
		}
	}
}

func (w *CleanupWorker) Stop() {
	log.Println("Cleanup worker stopping...")
}

func (w *CleanupWorker) cleanup() {
	log.Println("Running cleanup task...")
}