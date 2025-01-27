package pkg

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

type SearchCriteria struct {
	Name    string
	MinSize int64
	MaxSize int64
	After   time.Time
	Before  time.Time
}

func SearchFiles(root string, criteria SearchCriteria) ([]string, error) {
	resultChan := make(chan string)
	errChan := make(chan error, 1)
	var wg sync.WaitGroup
	var results []string

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				files, err := os.ReadDir(path)
				if err != nil {
					return err
				}

				for _, file := range files {
					if file.IsDir() {
						continue
					}

					// Get detailed file info for size and modification time
					fileInfo, err := file.Info()
					if err != nil {
						continue
					}

					// Check if file matches all specified criteria
					matches := true

					// Name check
					if criteria.Name != "" && file.Name() != criteria.Name {
						matches = false
					}

					// Size check
					if criteria.MinSize > 0 && fileInfo.Size() < criteria.MinSize {
						matches = false
					}
					if criteria.MaxSize > 0 && fileInfo.Size() > criteria.MaxSize {
						matches = false
					}

					// Modification time check
					modTime := fileInfo.ModTime()
					if !criteria.After.IsZero() && modTime.Before(criteria.After) {
						matches = false
					}
					if !criteria.Before.IsZero() && modTime.After(criteria.Before) {
						matches = false
					}

					if matches {
						select {
						case resultChan <- filepath.Join(path, file.Name()):
						case <-errChan:
							return filepath.SkipAll
						}
					}
				}
			}
			return nil
		})
		if err != nil {
			select {
			case errChan <- err:
			default:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(resultChan)
		close(errChan)
	}()

	var err error
	for {
		select {
		case path, ok := <-resultChan:
			if !ok {
				return results, err
			}
			results = append(results, path)
		case e, ok := <-errChan:
			if ok {
				err = e
			}
		}
	}
}
