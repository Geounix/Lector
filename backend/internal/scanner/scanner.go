package scanner

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type ComicParser struct {
	volumeRegex   *regexp.Regexp
	chapterRegex  *regexp.Regexp
	specialRegex  *regexp.Regexp
}

type ComicInfo struct {
	SeriesName string
	Volume     *int
	Chapter    *int
	Special    *string
	Title      string
	Hash       string
	Format     string
}

func NewComicParser() *ComicParser {
	return &ComicParser{
		volumeRegex:  regexp.MustCompile(`(?i)vol\.?\s*(\d+)`),
		chapterRegex: regexp.MustCompile(`(?i)ch\.?\s*(\d+)`),
		specialRegex: regexp.MustCompile(`(?i)(one-shot|special|extra|omake)`),
	}
}

func (p *ComicParser) ParseFilename(filename string) *ComicInfo {
	basename := strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))

	info := &ComicInfo{
		Title: basename,
	}

	if matches := p.volumeRegex.FindStringSubmatch(basename); len(matches) > 1 {
		vol := regexp.MustCompile(`\d+`).FindString(matches[1])
		if v, err := parseInt(vol); err == nil {
			info.Volume = &v
		}
	}

	if matches := p.chapterRegex.FindStringSubmatch(basename); len(matches) > 1 {
		ch := regexp.MustCompile(`\d+`).FindString(matches[1])
		if c, err := parseInt(ch); err == nil {
			info.Chapter = &c
		}
	}

	if matches := p.specialRegex.FindStringSubmatch(basename); len(matches) > 1 {
		info.Special = &matches[1]
	}

	info.SeriesName = p.extractSeriesName(basename)
	info.Format = strings.ToLower(filepath.Ext(filename)[1:])

	return info
}

func (p *ComicParser) extractSeriesName(basename string) string {
	patterns := []string{
		`^(.+?)\s+(?:vol\.?\s*\d+|ch\.?\s*\d+)`,
		`^(.+?)\s+\d+$`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(basename); len(matches) > 1 {
			return strings.TrimSpace(matches[1])
		}
	}

	return basename
}

func CalculateHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func ExtractPageCount(filePath string) (int, error) {
	ext := strings.ToLower(filepath.Ext(filePath))
	if ext == ".cbz" || ext == ".zip" {
		return countZipPages(filePath)
	}
	return 0, nil
}

func countZipPages(filePath string) (int, error) {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return 0, err
	}
	defer reader.Close()

	count := 0
	for _, file := range reader.File {
		if isImageFile(file.Name) {
			count++
		}
	}

	return count, nil
}

func isImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif"
}

func parseInt(s string) (int, error) {
	var result int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			result = result*10 + int(c-'0')
		} else {
			return 0, nil
		}
	}
	return result, nil
}