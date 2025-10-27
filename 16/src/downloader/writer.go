package downloader

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func (d *Downloader) getLocalPath(urlStr string) string {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		hash := fmt.Sprintf("%x", md5.Sum([]byte(urlStr)))
		return filepath.Join(d.config.OutputDir, hash+".html")
	}

	cleanPath := parsedURL.Path
	if cleanPath == "" || cleanPath == "/" {
		cleanPath = "/index.html"
	}

	localPath := filepath.Join(d.config.OutputDir, parsedURL.Host, cleanPath)

	if strings.HasSuffix(cleanPath, "/") {
		localPath = filepath.Join(localPath, "index.html")
	} else if filepath.Ext(cleanPath) == "" && !d.isResource(cleanPath) {
		localPath += ".html"
	}

	return localPath
}

func (d *Downloader) isResource(path string) bool {
	path = strings.ToLower(path)

	// Common resource paths
	patterns := []string{"/img/", "/css/", "/js/", "/assets/", "/static/", "/media/"}
	for _, pattern := range patterns {
		if strings.Contains(path, pattern) {
			return true
		}
	}

	// Known file extensions
	ext := strings.ToLower(filepath.Ext(path))
	extensions := []string{
		".css", ".js", ".png", ".jpg", ".gif", ".svg", ".ico",
		".woff", ".ttf", ".pdf", ".zip", ".xml", ".json",
	}

	for _, e := range extensions {
		if ext == e {
			return true
		}
	}

	return false
}

func (d *Downloader) getRelativePath(fromPath, toPath string) string {
	relPath, err := filepath.Rel(filepath.Dir(fromPath), toPath)
	if err != nil {
		return ""
	}

	return strings.ReplaceAll(relPath, "\\", "/")
}

func (d *Downloader) downloadFile(urlStr string) (string, []byte, error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", nil, fmt.Errorf("Unable to create request: %w", err)
	}

	req.Header.Set("User-Agent", "wget/1.0")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", nil, fmt.Errorf("Unable to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", nil, fmt.Errorf("Got status code: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil, fmt.Errorf("Unable to read response: %w", err)
	}

	localPath := d.getLocalPath(urlStr)

	dir := filepath.Dir(localPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", nil, fmt.Errorf("Unable to create directory %s: %w", dir, err)
	}

	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return "", nil, fmt.Errorf("Unable to save file %s: %w", localPath, err)
	}

	d.incrementDownloadedCount()
	return localPath, content, nil
}

func (d *Downloader) saveUpdatedContent(localPath string, content []byte) error {
	if err := os.WriteFile(localPath, content, 0644); err != nil {
		return fmt.Errorf("Unable to save updated file %s: %w", localPath, err)
	}
	return nil
}