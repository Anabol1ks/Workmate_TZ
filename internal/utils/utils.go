package utils

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var allowedExtensions = []string{".pdf", ".jpeg", ".jpg"}

func getProjectTempDir() string {
	dir := "temp"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		_ = os.Mkdir(dir, 0755)
	}
	return dir
}

func IsAllowedExtension(url string) bool {
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(strings.ToLower(url), ext) {
			return true
		}
	}
	return false
}

func DownloadFile(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("file not available: " + resp.Status)
	}

	parts := strings.Split(url, "/")
	name := parts[len(parts)-1]
	localPath := filepath.Join(getProjectTempDir(), name)
	out, err := os.Create(localPath)
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}
	return localPath, nil
}

func CreateZipArchive(files []string, archiveName string) (string, error) {
	archivePath := filepath.Join(getProjectTempDir(), archiveName)
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer zipFile.Close()
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		f, err := os.Open(file)
		if err != nil {
			return "", err
		}
		defer f.Close()
		info, err := f.Stat()
		if err != nil {
			return "", err
		}
		w, err := zipWriter.Create(info.Name())
		if err != nil {
			return "", err
		}
		_, err = io.Copy(w, f)
		if err != nil {
			return "", err
		}
	}
	return archivePath, nil
}
