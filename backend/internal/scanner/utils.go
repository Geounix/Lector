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

func CalculateHashPartial(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	data := make([]byte, 1024)
	n, err := file.Read(data)
	if err != nil && err != io.EOF {
		return "", err
	}

	h := sha256.New()
	h.Write(data[:n])
	return hex.EncodeToString(h.Sum(nil)), nil
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

func IsImageFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".webp" || ext == ".gif"
}

func GetFileSize(path string) int64 {
	info, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return info.Size()
}

func GenerateSortTitle(title string) string {
	sortTitle := strings.ToLower(title)
	sortTitle = strings.ReplaceAll(sortTitle, "the ", "")
	sortTitle = strings.ReplaceAll(sortTitle, "a ", "")
	sortTitle = strings.ReplaceAll(sortTitle, "an ", "")
	sortTitle = strings.TrimSpace(sortTitle)
	return sortTitle
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

var volumeRegex = regexp.MustCompile(`(?i)vol\.?\s*(\d+)`)
var chapterRegex = regexp.MustCompile(`(?i)ch\.?\s*(\d+)`)
var specialRegex = regexp.MustCompile(`(?i)(one-shot|special|extra|omake)`)

func ExtractVolumeNumber(filename string) int {
	matches := volumeRegex.FindStringSubmatch(filename)
	if len(matches) > 1 {
		digits := regexp.MustCompile(`\d+`).FindString(matches[1])
		if v, err := parseInt(digits); err == nil {
			return v
		}
	}
	return 1
}

func ExtractChapterNumber(filename string) int {
	matches := chapterRegex.FindStringSubmatch(filename)
	if len(matches) > 1 {
		digits := regexp.MustCompile(`\d+`).FindString(matches[1])
		if c, err := parseInt(digits); err == nil {
			return c
		}
	}
	return 0
}