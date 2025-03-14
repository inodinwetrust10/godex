package version

import (
	"bufio"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CreateFile(filePath, versionID, message string) (VersionMetaData, error) {
	fileDir, err := GetVersionPath(filePath)
	if err != nil {
		return VersionMetaData{}, err
	}
	lastFilePath := ReturnLastFilePath(fileDir)
	if lastFilePath == "" {
		return VersionMetaData{}, err
	}
	isRequired, err := checkDiffs(filePath, lastFilePath)
	if err != nil {
		return VersionMetaData{}, err
	}
	// if a version without change already exists it return an error
	if isRequired == false {
		return VersionMetaData{}, errors.New("A version already exists")
	}
	meta, err := saveFile(filePath, versionID, message, fileDir)
	return meta, err
}

// /////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////
func ListAllVersions(fileDir string) (*[]VersionMetaData, error) {
	metaDataFilePath := filepath.Join(fileDir, "version.json")

	var listAllVersions []VersionMetaData

	data, err := os.ReadFile(metaDataFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"No version currently exists of this file",
			)
		}
		return nil, err
	}
	err = json.Unmarshal(data, &listAllVersions)
	if err != nil {
		return nil, err
	}
	return &listAllVersions, nil
}

// /////////////////////////////////////////////////////////////////////////////////
// /////////////////////////////////////////////////////////////////////////////////

func FileDiff(path1, path2 string) (DiffResult, error) {
	result := DiffResult{
		Identical: true,
		DiffLines: []LineDiff{},
	}

	info1, err := os.Stat(path1)
	if err != nil {
		return result, fmt.Errorf("error accessing first file: %w", err)
	}

	_, err = os.Stat(path2)
	if err != nil {
		return result, fmt.Errorf("error accessing second file: %w", err)
	}

	if info1.Size() < 10*1024*1024 {
		return compareFilesDetailed(path1, path2)
	}

	checksum1, err := calculateMD5(path1)
	if err != nil {
		return result, err
	}

	checksum2, err := calculateMD5(path2)
	if err != nil {
		return result, err
	}

	if checksum1 != checksum2 {
		result.Identical = false
		result.DiffType = "content"
		result.Message = "Files have different content (detected by checksum)"

		detailedResult, err := compareFilesDetailed(path1, path2)
		if err != nil {
			return result, err
		}
		result.DiffLines = detailedResult.DiffLines
		return result, nil
	}

	result.Message = "Files are identical"
	return result, nil
}

// calculateMD5 calculates the MD5 checksum of a file
func calculateMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func compareFilesDetailed(path1, path2 string) (DiffResult, error) {
	result := DiffResult{
		Identical: true,
		DiffLines: []LineDiff{},
	}

	file1, err := os.Open(path1)
	if err != nil {
		return result, err
	}
	defer file1.Close()

	file2, err := os.Open(path2)
	if err != nil {
		return result, err
	}
	defer file2.Close()

	scanner1 := bufio.NewScanner(file1)
	scanner2 := bufio.NewScanner(file2)

	lineNum := 0

	for {
		hasLine1 := scanner1.Scan()
		hasLine2 := scanner2.Scan()
		lineNum++

		if !hasLine1 && !hasLine2 {
			break
		}

		if hasLine1 != hasLine2 {
			result.Identical = false
			result.DiffType = "line"

			line1 := ""
			line2 := ""

			if hasLine1 {
				line1 = scanner1.Text()
			}

			if hasLine2 {
				line2 = scanner2.Text()
			}

			result.DiffLines = append(result.DiffLines, LineDiff{
				LineNumber: lineNum,
				Line1:      line1,
				Line2:      line2,
			})

			continue
		}

		// Both files have lines at this point, so compare them
		line1 := scanner1.Text()
		line2 := scanner2.Text()

		if line1 != line2 {
			result.Identical = false
			result.DiffType = "line"

			result.DiffLines = append(result.DiffLines, LineDiff{
				LineNumber: lineNum,
				Line1:      line1,
				Line2:      line2,
			})
		}
	}

	// Check for scanner errors
	if err := scanner1.Err(); err != nil {
		return result, err
	}
	if err := scanner2.Err(); err != nil {
		return result, err
	}

	if !result.Identical {
		result.Message = fmt.Sprintf("Found %d different lines", len(result.DiffLines))
	} else {
		result.Message = "Files are identical"
	}

	return result, nil
}

func FormatDiffResult(result DiffResult) string {
	var sb strings.Builder

	if result.Identical {
		sb.WriteString("Files are identical\n")
		return sb.String()
	}

	sb.WriteString(result.Message + "\n\n")

	if len(result.DiffLines) > 0 {
		sb.WriteString("Differences found:\n")

		for _, diff := range result.DiffLines {
			sb.WriteString(fmt.Sprintf("Line %d:\n", diff.LineNumber))
			sb.WriteString(fmt.Sprintf("  File 1: %s\n", diff.Line1))
			sb.WriteString(fmt.Sprintf("  File 2: %s\n\n", diff.Line2))
		}
	}

	return sb.String()
}

// ////////////////////////////////////////////////////////////////////////////////////////////////
// ////////////////////////////////////////////////////////////////////////////////////////////////
func ClearAllVersion(dirPath string) error {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory: %w", err)
	}
	deletedCount := 0
	errorCount := 0

	for _, entry := range entries {
		if entry.IsDir() {
			fmt.Printf("Skipping subdirectory: %s\n", entry.Name())
			continue
		}
		filePath := filepath.Join(dirPath, entry.Name())

		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Error deleting %s: %v\n", filePath, err)
			errorCount++
		} else {
			fmt.Printf("Deleted: %s\n", filePath)
			deletedCount++
		}
	}
	fmt.Printf("\nSummary: Deleted %d files with %d errors\n", deletedCount, errorCount)
	return nil
}

// ///////////////////////////////////////////////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////////////////////////////////////////////
func ClearVersion(dirPath, versionID string) error {
	versionPath := filepath.Join(dirPath, versionID)

	if _, err := os.Stat(versionPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("no file exists with versionID %s", versionID)
		}
		return fmt.Errorf("failed to access file %s: %v", versionID, err)
	}

	if err := os.Remove(versionPath); err != nil {
		return fmt.Errorf("failed to remove file %s: %v", versionID, err)
	}

	allVersions, err := ListAllVersions(dirPath)
	if err != nil {
		return fmt.Errorf("failed to list versions: %v", err)
	}

	if allVersions == nil {
		return fmt.Errorf("version list is nil")
	}

	i := 0
	for _, version := range *allVersions {
		if version.ID != versionID {
			(*allVersions)[i] = version
			i++
		}
	}
	*allVersions = (*allVersions)[:i]

	jsonPath := filepath.Join(dirPath, "version.json")
	jsonData, err := json.MarshalIndent(*allVersions, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %v", err)
	}

	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata to file %s: %v", jsonPath, err)
	}

	return nil
}
