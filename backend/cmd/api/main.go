package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lector-comics/lector/internal/config"
	"github.com/lector-comics/lector/internal/handlers"
	"github.com/lector-comics/lector/internal/middleware"
	"github.com/lector-comics/lector/internal/repository"
	"github.com/lector-comics/lector/internal/services"
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
	svc := services.NewService(repo)
	h := handlers.NewHandler(svc, cfg.JWTSecret)

	app := fiber.New(fiber.Config{
		AppName: "Lector Comics API v0.1",
	})

	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	app.Get("/health", h.HealthCheck)

	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", h.Register)
	auth.Post("/login", h.Login)
	auth.Post("/logout", h.Logout)
	auth.Post("/refresh", h.RefreshToken)

	authenticated := api.Group("", middleware.JWTAuth(cfg.JWTSecret), func(c *fiber.Ctx) error {
		c.Locals("jwt_secret", cfg.JWTSecret)
		return c.Next()
	})

	library := authenticated.Group("/library")
	library.Get("/", h.GetLibraries)
	library.Post("/", h.CreateLibrary)
	library.Patch("/:id", h.UpdateLibrary)
	library.Delete("/:id", h.DeleteLibrary)

	series := authenticated.Group("/series")
	series.Get("/", h.GetSeries)
	series.Get("/:id", h.GetSeriesByID)
	series.Post("/scan", h.ScanLibrary)

	reader := authenticated.Group("/reader")
	reader.Get("/chapter/:id", h.GetChapter)
	reader.Get("/chapter/:id/page/:page", h.GetPage)
	reader.Post("/progress", h.UpdateProgress)

	search := authenticated.Group("/search")
	search.Get("/", h.Search)

	metadata := authenticated.Group("/metadata")
	metadata.Get("/:seriesId", h.GetMetadata)
	metadata.Post("/refresh/:seriesId", h.RefreshMetadata)

	stats := authenticated.Group("/stats")
	stats.Get("/", h.GetStats)

	user := authenticated.Group("/user")
	user.Get("/", h.GetCurrentUser)
	user.Patch("/", h.UpdateUser)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Listen(":" + cfg.Port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")
	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}
}