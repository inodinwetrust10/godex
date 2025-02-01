package pkg

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type ProgressReader struct {
	io.Reader
	Total     int64
	ReadBytes int64
}

func (pr *ProgressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.ReadBytes += int64(n)
	return n, err
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func UploadFile(filePath string) error {
	ctx := context.Background()
	configPath := filepath.Join(GetConfigDir(), "credentials.json")

	b, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading credentials.json: %v", err)
	}

	config, err := google.ConfigFromJSON(b, drive.DriveFileScope)
	if err != nil {
		return fmt.Errorf("error parsing credentials.json: %v", err)
	}

	client := getClient(config)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("error creating Drive service: %v", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("error getting file info: %v", err)
	}

	pr := &ProgressReader{
		Reader: file,
		Total:  fileInfo.Size(),
	}

	done := make(chan bool)
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				percent := float64(pr.ReadBytes) / float64(pr.Total) * 100
				fmt.Printf("\rUploading... %s/%s (%.2f%%)",
					formatBytes(pr.ReadBytes),
					formatBytes(pr.Total),
					percent)
			case <-done:
				fmt.Printf("\rUpload complete! %s uploaded\n", formatBytes(pr.Total))
				return
			}
		}
	}()

	driveFile := &drive.File{Name: fileInfo.Name()}

	_, err = srv.Files.Create(driveFile).Media(pr).Context(ctx).Do()

	done <- true
	if err != nil {
		return fmt.Errorf("error uploading file: %v", err)
	}

	return nil
}
