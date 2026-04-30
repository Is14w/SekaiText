package service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Downloader handles HTTP downloads from CDN.
type Downloader struct {
	client  *http.Client
	dataDir string
}

// ProgressTracker tracks download progress.
type ProgressTracker struct {
	mu       sync.RWMutex
	current  int
	total    int
	message  string
	done     bool
}

// NewDownloader creates a new Downloader.
func NewDownloader(dataDir string) *Downloader {
	return &Downloader{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		dataDir: dataDir,
	}
}

// NewProgressTracker creates a new ProgressTracker.
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{}
}

// SetTotal sets the total number of steps.
func (pt *ProgressTracker) SetTotal(total int) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.total = total
	pt.current = 0
	pt.done = false
}

// Advance increments the progress counter.
func (pt *ProgressTracker) Advance(msg string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.current++
	pt.message = msg
}

// Done marks the progress as complete.
func (pt *ProgressTracker) Done() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	pt.done = true
	pt.message = "完成"
}

// Status returns the current progress status.
func (pt *ProgressTracker) Status() (int, int, string, bool) {
	pt.mu.RLock()
	defer pt.mu.RUnlock()
	return pt.current, pt.total, pt.message, pt.done
}

// DownloadJSON downloads a JSON file from URL and saves to dataDir.
// Returns file path.
func (d *Downloader) DownloadJSON(url, fileName string) (string, error) {
	filePath := filepath.Join(d.dataDir, fileName)

	// Check if already cached
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Cache hit: %s", fileName)
		return filePath, nil
	}

	log.Printf("Downloading: %s", url)
	resp, err := d.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Downloaded: %s", fileName)
	return filePath, nil
}

// DownloadProgressFn is called with (bytesRead, totalBytes) during download.
type DownloadProgressFn func(read, total int64)

// DownloadJSONToDir downloads a JSON file from URL to a specific directory.
// If progress is non-nil, it is called periodically with bytes read and total bytes.
func (d *Downloader) DownloadJSONToDir(url, dir, fileName string, progress DownloadProgressFn) (string, error) {
	filePath := filepath.Join(dir, fileName)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Check if already cached
	if _, err := os.Stat(filePath); err == nil {
		log.Printf("Cache hit: %s", filePath)
		if progress != nil {
			// Get file size for progress
			info, _ := os.Stat(filePath)
			progress(info.Size(), info.Size())
		}
		return filePath, nil
	}

	log.Printf("Downloading: %s -> %s", url, filePath)
	resp, err := d.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download returned status %d", resp.StatusCode)
	}

	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	var reader io.Reader = resp.Body
	if progress != nil && resp.ContentLength > 0 {
		reader = &progressReader{reader: resp.Body, total: resp.ContentLength, fn: progress}
	}

	if _, err := io.Copy(out, reader); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	if progress != nil {
		progress(resp.ContentLength, resp.ContentLength)
	}

	log.Printf("Downloaded: %s", filePath)
	return filePath, nil
}

type progressReader struct {
	reader io.Reader
	total  int64
	read   int64
	fn     DownloadProgressFn
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.reader.Read(b)
	p.read += int64(n)
	p.fn(p.read, p.total)
	return n, err
}

// DownloadAndParseJSON downloads a JSON file and parses it.
func (d *Downloader) DownloadAndParseJSON(url, fileName string, target interface{}) error {
	filePath, err := d.DownloadJSON(url, fileName)
	if err != nil {
		return err
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read downloaded file: %w", err)
	}

	return json.Unmarshal(data, target)
}

// UpdateAll performs a full metadata update from CDN.
func (lm *ListManager) UpdateAll(catalogDir string, pt *ProgressTracker) {
	lm.UpdateAllFromCDN(catalogDir, pt)
}
