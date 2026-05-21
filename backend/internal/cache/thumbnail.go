package cache

import (
	"archive/zip"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/redis/go-redis/v9"
)

type ThumbnailService struct {
	redis  *redis.Client
	cacheDir string
}

func NewThumbnailService(redis *redis.Client, cacheDir string) *ThumbnailService {
	return &ThumbnailService{
		redis:    redis,
		cacheDir: cacheDir,
	}
}

type ThumbnailConfig struct {
	Width  int
	Height int
}

func (t *ThumbnailService) DefaultConfig() ThumbnailConfig {
	return ThumbnailConfig{Width: 512, Height: 768}
}

func (t *ThumbnailService) GenerateThumbnail(comicPath string, seriesID int64) (string, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("thumbnail:%d", seriesID)

	cached, err := t.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		return cached, nil
	}

	ext := strings.ToLower(filepath.Ext(comicPath))
	if ext != ".cbz" && ext != ".zip" {
		return "", fmt.Errorf("unsupported format: %s", ext)
	}

	reader, err := zip.OpenReader(comicPath)
	if err != nil {
		return "", fmt.Errorf("failed to open archive: %w", err)
	}
	defer reader.Close()

	var firstImage string
	for _, file := range reader.File {
		if isImageFile(file.Name) && !strings.HasPrefix(filepath.Base(file.Name), ".") {
			firstImage = file.Name
			break
		}
	}

	if firstImage == "" {
		return "", fmt.Errorf("no images found in archive")
	}

	for _, file := range reader.File {
		if file.Name == firstImage {
			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open image: %w", err)
			}

			img, format, err := image.Decode(rc)
			rc.Close()
			if err != nil {
				return "", fmt.Errorf("failed to decode image: %w", err)
			}

			outputPath, err := t.saveThumbnail(img, seriesID, format)
			if err != nil {
				return "", fmt.Errorf("failed to save thumbnail: %w", err)
			}

			t.redis.Set(ctx, cacheKey, outputPath, 30*24*60*60*1000000000)
			return outputPath, nil
		}
	}

	return "", fmt.Errorf("image not found in archive")
}

func (t *ThumbnailService) saveThumbnail(img image.Image, seriesID int64, format string) (string, error) {
	if err := os.MkdirAll(t.cacheDir, 0755); err != nil {
		return "", err
	}

	outputPath := filepath.Join(t.cacheDir, fmt.Sprintf("series_%d.jpg", seriesID))

	var file *os.File
	var err error

	switch format {
	case "jpeg", "jpg":
		file, err = os.Create(outputPath)
		if err != nil {
			return "", err
		}
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	case "png":
		file, err = os.Create(outputPath)
		if err != nil {
			return "", err
		}
		err = png.Encode(file, img)
	default:
		file, err = os.Create(outputPath)
		if err != nil {
			return "", err
		}
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	}

	if err != nil {
		file.Close()
		return "", err
	}

	file.Close()
	return outputPath, nil
}

func (t *ThumbnailService) DeleteThumbnail(seriesID int64) error {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("thumbnail:%d", seriesID)

	t.redis.Del(ctx, cacheKey)

	outputPath := filepath.Join(t.cacheDir, fmt.Sprintf("series_%d.jpg", seriesID))
	os.Remove(outputPath)

	return nil
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif"
}