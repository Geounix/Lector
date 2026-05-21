package reader

import (
	"archive/zip"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type ComicReader struct {
	cache map[string][]string
}

func NewComicReader() *ComicReader {
	return &ComicReader{
		cache: make(map[string][]string),
	}
}

func (r *ComicReader) GetPages(filePath string) ([]string, error) {
	if pages, ok := r.cache[filePath]; ok {
		return pages, nil
	}

	ext := filepath.Ext(filePath)
	switch ext {
	case ".cbz", ".zip":
		return r.getZipPages(filePath)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}

func (r *ComicReader) getZipPages(filePath string) ([]string, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	var pages []string
	for _, file := range reader.File {
		if isImageFile(file.Name) {
			pages = append(pages, file.Name)
		}
	}

	r.cache[filePath] = pages
	return pages, nil
}

func (r *ComicReader) GetPage(filePath string, pageIndex int) ([]byte, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	pages, err := r.GetPages(filePath)
	if err != nil {
		return nil, err
	}

	if pageIndex < 0 || pageIndex >= len(pages) {
		return nil, fmt.Errorf("page index out of range")
	}

	for _, file := range reader.File {
		if file.Name == pages[pageIndex] {
			return extractImageData(file)
		}
	}

	return nil, fmt.Errorf("page not found")
}

func extractImageData(file *zip.File) ([]byte, error) {
	rc, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif"
}

func SaveImage(img image.Image, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	case ".png":
		return png.Encode(file, img)
	default:
		return fmt.Errorf("unsupported image format: %s", ext)
	}
}