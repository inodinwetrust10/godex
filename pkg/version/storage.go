package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// //////////////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////////////
func saveFile(
	filePath,
	versionID,
	message,
	versionPathDir string,
) (VersionMetaData, error) {
	versionFilePath := filepath.Join(versionPathDir, versionID)

	sourceFile, err := os.Open(filePath)
	if err != nil {
		return VersionMetaData{}, fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(versionFilePath)
	if err != nil {
		return VersionMetaData{}, fmt.Errorf("failed to create version file: %w", err)
	}
	defer destinationFile.Close()

	hasher := sha256.New()

	tee := io.TeeReader(sourceFile, hasher)

	size, err := io.Copy(destinationFile, tee)
	if err != nil {
		return VersionMetaData{}, fmt.Errorf("failed to copy file: %w", err)
	}

	checksum := hex.EncodeToString(hasher.Sum(nil))

	metadata := VersionMetaData{
		ID:        versionID,
		CreatedAt: time.Now(),
		Message:   message,
		Size:      int64(size),
		Checksum:  checksum,
	}

	if err = saveMetaData(versionPathDir, metadata); err != nil {
		return VersionMetaData{}, fmt.Errorf("failed to update meta data: %w", err)
	}

	if err = updateGlobalIndex(versionID, filePath); err != nil {
		return VersionMetaData{}, fmt.Errorf("unable to update version index: %w", err)
	}

	return metadata, nil
}

// //////////////////////////////////////////////////////////////////////////////////////
// //////////////////////////////////////////////////////////////////////////////////////

func saveMetaData(filePath string, metadata VersionMetaData) error {
	fullPath := filepath.Join(filePath, "version.json")
	var metadataEntries []VersionMetaData

	if _, err := os.Stat(fullPath); err == nil {
		fileData, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read existing metadata file: %w", err)
		}

		// Unmarshal existing data
		if err := json.Unmarshal(fileData, &metadataEntries); err != nil {
			var singleMetadata VersionMetaData
			if err := json.Unmarshal(fileData, &singleMetadata); err != nil {
				return fmt.Errorf("failed to parse existing metadata: %w", err)
			}
			metadataEntries = append(metadataEntries, singleMetadata)
		}
		metadataEntries = append(metadataEntries, metadata)
	} else {
		// File doesnt exist start with just the new metadata
		metadataEntries = []VersionMetaData{metadata}
	}

	jsonData, err := json.MarshalIndent(metadataEntries, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata to JSON: %w", err)
	}

	if err := os.WriteFile(fullPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata to file %s: %w", fullPath, err)
	}

	return nil
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////

func updateGlobalIndex(versionID, filePath string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	dirPath := filepath.Join(homeDir, ".config", "godex")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	globalIndexPath := filepath.Join(dirPath, "global.json")
	indexMap := make(map[string]*GlobalIndex)

	if _, err := os.Stat(globalIndexPath); err == nil {
		fileData, err := os.ReadFile(globalIndexPath)
		if err != nil {
			return fmt.Errorf("failed to read global index file: %w", err)
		}

		var indices []GlobalIndex
		if err := json.Unmarshal(fileData, &indices); err != nil {
			var singleIndex GlobalIndex
			if err := json.Unmarshal(fileData, &singleIndex); err != nil {
				return fmt.Errorf("failed to parse global index: %w", err)
			}
			indices = []GlobalIndex{singleIndex}
		}

		for i := range indices {
			indexMap[indices[i].OriginalFilePath] = &indices[i]
		}
	}

	now := time.Now()
	if idx, exists := indexMap[filePath]; exists {
		versionExists := false
		for _, v := range idx.Versions {
			if v == versionID {
				versionExists = true
				break
			}
		}

		if !versionExists {
			idx.Versions = append(idx.Versions, versionID)
		}
		idx.LastUpdatedAt = now
	} else {
		indexMap[filePath] = &GlobalIndex{
			OriginalFilePath: filePath,
			Versions:         []string{versionID},
			LastUpdatedAt:    now,
		}
	}

	var updatedIndices []GlobalIndex
	for _, idx := range indexMap {
		updatedIndices = append(updatedIndices, *idx)
	}

	jsonData, err := json.MarshalIndent(updatedIndices, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal global index to JSON: %w", err)
	}

	if err := os.WriteFile(globalIndexPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write global index to file: %w", err)
	}

	return nil
}
