package downloader

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"wb_l2/16/src/config"
)

const (
	maxQueueSize = 10000
	pollInterval = 100 * time.Millisecond
)

type downloadTask struct {
	url   string
	depth int
}

type Downloader struct {
	config          config.Config
	client          *http.Client
	visited         map[string]bool
	visitedMutex    sync.RWMutex
	baseURL         *url.URL
	downloadedCount int64
	taskChan        chan downloadTask
	workerWg        sync.WaitGroup
	activeTasks     int64
}

func NewDownloader(config config.Config) *Downloader {
	return &Downloader{
		config:   config,
		client:   &http.Client{Timeout: config.Timeout},
		visited:  make(map[string]bool),
		taskChan: make(chan downloadTask, maxQueueSize),
	}
}

func (d *Downloader) isVisited(urlStr string) bool {
	d.visitedMutex.RLock()
	defer d.visitedMutex.RUnlock()
	return d.visited[urlStr]
}

func (d *Downloader) markVisited(urlStr string) {
	d.visitedMutex.Lock()
	defer d.visitedMutex.Unlock()
	d.visited[urlStr] = true
}

func (d *Downloader) isSameDomain(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return parsedURL.Host == d.baseURL.Host
}

func (d *Downloader) GetDownloadedCount() int {
	return int(atomic.LoadInt64(&d.downloadedCount))
}

func (d *Downloader) incrementDownloadedCount() {
	atomic.AddInt64(&d.downloadedCount, 1)
}

func (d *Downloader) Download(urlStr string) error {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	d.baseURL = parsedURL

	if err := os.MkdirAll(d.config.OutputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	d.startWorkers()

	d.addTask(urlStr, 0)

	d.waitForCompletion()

	close(d.taskChan)
	d.workerWg.Wait()

	return nil
}

func (d *Downloader) startWorkers() {
	for i := 0; i < d.config.MaxConcurrent; i++ {
		d.workerWg.Add(1)
		go d.worker()
	}
}

func (d *Downloader) worker() {
	defer d.workerWg.Done()

	for task := range d.taskChan {
		d.processURL(task.url, task.depth)
	}
}

func (d *Downloader) addTask(urlStr string, depth int) {
	if depth > d.config.MaxDepth || d.isVisited(urlStr) {
		return
	}

	d.markVisited(urlStr)

	atomic.AddInt64(&d.activeTasks, 1)

	task := downloadTask{url: urlStr, depth: depth}

	d.taskChan <- task
}

func (d *Downloader) processURL(urlStr string, depth int) {
	defer atomic.AddInt64(&d.activeTasks, -1)

	if d.config.Verbose {
		fmt.Printf("Downloading: %s, depth = %d\n", urlStr, depth)
	}

	localPath, content, err := d.downloadFile(urlStr)
	if err != nil {
		if d.config.Verbose {
			fmt.Printf("Error: %s: %v\n", urlStr, err)
		}

		return
	}

	if d.config.Verbose {
		fmt.Printf("Saved file: %s\n", localPath)
	}

	contentType := http.DetectContentType(content)
	if strings.Contains(contentType, "text/html") {
		links := d.extractLinks(content, urlStr)

		updatedContent := d.updateLinks(content, urlStr)
		if err := d.saveUpdatedContent(localPath, updatedContent); err != nil && d.config.Verbose {
			fmt.Printf("Error updating links in %s: %v\n", localPath, err)
		}

		for _, link := range links {
			if d.isSameDomain(link) {
				d.addTask(link, depth+1)
			}
		}
	}
}

func (d *Downloader) waitForCompletion() {
	for {
		active := atomic.LoadInt64(&d.activeTasks)
		if active == 0 {
			break
		}

		time.Sleep(pollInterval)
	}
}
