package workers

import (
	"log"

	"github.com/lector-comics/lector/internal/repository"
)

type ThumbnailWorker struct {
	repo *repository.Repository
}

func NewThumbnailWorker(repo *repository.Repository) *ThumbnailWorker {
	return &ThumbnailWorker{repo: repo}
}

func (w *ThumbnailWorker) Start() {
	log.Println("Thumbnail worker started")
}

func (w *ThumbnailWorker) Stop() {
	log.Println("Thumbnail worker stopping...")
}