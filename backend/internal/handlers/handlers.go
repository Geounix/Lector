package handlers

import (
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lector-comics/lector/internal/middleware"
	"github.com/lector-comics/lector/internal/reader"
	"github.com/lector-comics/lector/internal/services"
)

type Handler struct {
	svc       *services.Service
	jwtSecret string
	comicReader *reader.ComicReader
}

func NewHandler(svc *services.Service, jwtSecret string) *Handler {
	return &Handler{svc: svc, jwtSecret: jwtSecret, comicReader: reader.NewComicReader()}
}

func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "healthy",
		"service": "lector-comics-api",
		"version": "0.1.0",
	})
}

func (h *Handler) Register(c *fiber.Ctx) error {
	type req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if body.Username == "" || body.Email == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username, email and password are required",
		})
	}

	ctx := context.Background()
	user, err := h.svc.CreateUser(ctx, body.Username, body.Email, body.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, user.Role, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func (h *Handler) Login(c *fiber.Ctx) error {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	user, err := h.svc.GetUserByEmail(ctx, body.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to authenticate",
		})
	}

	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	if !h.svc.ValidatePassword(user, body.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	jwtSecret := h.jwtSecret
	token, err := middleware.GenerateToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}

	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate refresh token",
		})
	}

	return c.JSON(fiber.Map{
		"token":         token,
		"refresh_token": refreshToken,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"avatar":   user.Avatar,
			"theme":    user.Theme,
		},
	})
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"token": "new_token",
	})
}

func (h *Handler) GetLibraries(c *fiber.Ctx) error {
	ctx := context.Background()
	libraries, err := h.svc.GetLibraries(ctx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get libraries",
		})
	}

	if libraries == nil {
		libraries = []services.Library{}
	}

	return c.JSON(libraries)
}

func (h *Handler) CreateLibrary(c *fiber.Ctx) error {
	type req struct {
		Name string `json:"name"`
		Type string `json:"type"`
		Path string `json:"path"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx := context.Background()
	library, err := h.svc.CreateLibrary(ctx, body.Name, body.Type, body.Path)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create library",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(library)
}

func (h *Handler) UpdateLibrary(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "update library endpoint - todo"})
}

func (h *Handler) DeleteLibrary(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "delete library endpoint - todo"})
}

func (h *Handler) GetSeries(c *fiber.Ctx) error {
	ctx := context.Background()
	seriesList, err := h.svc.GetSeries(ctx, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get series",
		})
	}

	if seriesList == nil {
		seriesList = []services.Series{}
	}

	return c.JSON(seriesList)
}

func (h *Handler) GetSeriesByID(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "get series by id endpoint - todo"})
}

func (h *Handler) ScanLibrary(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "scan started",
		"status":  "queued",
	})
}

func (h *Handler) GetChapter(c *fiber.Ctx) error {
	chapterID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chapter ID",
		})
	}

	ctx := context.Background()
	chapter, err := h.svc.GetChapter(ctx, int64(chapterID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get chapter",
		})
	}

	if chapter == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chapter not found",
		})
	}

	userID := c.Locals("user_id").(int64)
	progress, _ := h.svc.GetReadingProgress(ctx, userID, int64(chapterID))

	return c.JSON(fiber.Map{
		"id":           chapter.ID,
		"volume_id":    chapter.VolumeID,
		"title":        chapter.Title,
		"file_path":    chapter.FilePath,
		"page_count":   chapter.PageCount,
		"current_page": progress.Page,
	})
}

func (h *Handler) GetPage(c *fiber.Ctx) error {
	chapterID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid chapter ID",
		})
	}

	pageStr := c.Params("page")
	pageIndex, err := strconv.Atoi(pageStr)
	if err != nil || pageIndex < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid page index",
		})
	}

	ctx := context.Background()
	chapter, err := h.svc.GetChapter(ctx, int64(chapterID))
	if err != nil || chapter == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Chapter not found",
		})
	}

	imageData, err := h.comicReader.GetPage(chapter.FilePath, pageIndex)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Page not found",
		})
	}

	c.Set("Content-Type", "image/jpeg")
	return c.Send(imageData)
}

func (h *Handler) UpdateProgress(c *fiber.Ctx) error {
	type req struct {
		ChapterID  int64  `json:"chapter_id"`
		Page       int    `json:"page"`
		Percentage float64 `json:"percentage"`
	}

	var body req
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	userID := c.Locals("user_id").(int64)
	ctx := context.Background()

	err := h.svc.UpsertReadingProgress(ctx, userID, body.ChapterID, body.Page, body.Percentage)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update progress",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Progress updated",
	})
}

func (h *Handler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	ctx := context.Background()
	results, err := h.svc.SearchSeries(ctx, query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Search failed",
		})
	}

	return c.JSON(fiber.Map{
		"results": results,
		"count":   len(results),
	})
}

func (h *Handler) GetMetadata(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "get metadata endpoint - todo"})
}

func (h *Handler) RefreshMetadata(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "refresh metadata endpoint - todo"})
}

func (h *Handler) GetStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	ctx := context.Background()

	stats, err := h.svc.GetStats(ctx, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get stats",
		})
	}

	return c.JSON(stats)
}

func (h *Handler) GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int64)
	ctx := context.Background()

	user, err := h.svc.GetUserByID(ctx, userID)
	if err != nil || user == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"id":       user.ID,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
		"avatar":   user.Avatar,
		"theme":    user.Theme,
		"language": user.Language,
	})
}

func (h *Handler) UpdateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"message": "update user endpoint - todo"})
}