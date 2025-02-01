package pkg

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func ZipFiles(outputFile string, files []string) error {
	zipFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err := addFileToZip(zipWriter, file); err != nil {
			return fmt.Errorf("failed to add %s to zip: %w", file, err)
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fmt.Printf("Adding file: %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info for %s: %w", filename, err)
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return fmt.Errorf("failed to create header for %s: %w", filename, err)
	}

	header.Name = filepath.Base(filename)
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return fmt.Errorf("failed to create header in zip for %s: %w", filename, err)
	}

	_, err = io.Copy(writer, file)
	if err != nil {
		return fmt.Errorf("failed to copy file %s to zip: %w", filename, err)
	}

	fmt.Printf("Successfully added file: %s\n", filename)
	return nil
}

func UnzipFile(inputFile, destination string) error {
	reader, err := zip.OpenReader(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		if err := extractFile(file, destination); err != nil {
			return fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}
	}

	return nil
}

func extractFile(file *zip.File, destination string) error {
	path := filepath.Join(destination, file.Name)

	if file.FileInfo().IsDir() {
		return os.MkdirAll(path, os.ModePerm)
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipFile, err := file.Open()
	if err != nil {
		return err
	}
	defer zipFile.Close()

	_, err = io.Copy(outFile, zipFile)
	return err
}

func ZipDirectory(outputFile string, directory string) error {
	zipFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	absDir, err := filepath.Abs(directory)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Ensure the input is a directory
	fileInfo, err := os.Stat(absDir)
	if err != nil {
		return fmt.Errorf("failed to stat directory: %w", err)
	}
	if !fileInfo.IsDir() {
		return fmt.Errorf("%s is not a directory", absDir)
	}

	baseDir := filepath.Dir(absDir) // Parent directory of the target directory

	err = filepath.Walk(absDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Get relative path from base directory to current file/directory
		relPath, err := filepath.Rel(baseDir, filePath)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create header: %w", err)
		}

		// Set header name to preserve directory structure
		header.Name = relPath

		// Handle directories by adding trailing slash and using Store method
		if info.IsDir() {
			header.Name += "/"
			header.Method = zip.Store
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create writer: %w", err)
		}

		if info.IsDir() {
			return nil // No content to write for directories
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		if err != nil {
			return fmt.Errorf("failed to copy file: %w", err)
		}

		fmt.Printf("Added: %s\n", relPath)
		return nil
	})
	if err != nil {
		return fmt.Errorf("error walking directory: %w", err)
	}

	return nil
}
