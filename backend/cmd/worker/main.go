package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lector-comics/lector/internal/config"
	"github.com/lector-comics/lector/internal/repository"
	"github.com/lector-comics/lector/internal/workers"
)

func main() {
	cfg := config.Load()

	db, err := repository.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redisClient := repository.NewRedisClient(cfg.RedisURL)
	defer redisClient.Close()

	repo := repository.NewRepository(db, redisClient)

	scannerWorker := workers.NewScannerWorker(repo)
	thumbnailWorker := workers.NewThumbnailWorker(repo)
	cleanupWorker := workers.NewCleanupWorker(repo)

	go scannerWorker.Start()
	go scannerWorker.SubscribeToScans()
	go thumbnailWorker.Start()
	go cleanupWorker.Start()

	log.Println("Worker service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down worker...")
}