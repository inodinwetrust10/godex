package version

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// ///////////////////////////////////////////////////////
// ///////////////////////////////////////////////////////
func getVersionPath(filePath string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}
	hash := sha256.Sum256([]byte(filePath))
	hashedName := hex.EncodeToString(hash[:])

	versionDir := filepath.Join(homeDir, ".config", "godex", "versions")
	versionFilePath := filepath.Join(versionDir, hashedName)

	err = os.MkdirAll(versionFilePath, 0755)
	if err != nil {
		return "", fmt.Errorf("could not create directory: %w", err)
	}

	return versionFilePath, nil
}

/////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////

func GenerateVersionID(filePath string) (string, error) {
	dirPath, err := getVersionPath(filePath)
	if err != nil {
		return "", err
	}

	versionFilePath := filepath.Join(dirPath, "version.json")

	nextVersion := 1

	if _, err := os.Stat(versionFilePath); err == nil {
		fileData, err := os.ReadFile(versionFilePath)
		if err != nil {
			return "", fmt.Errorf("failed to read version file: %w", err)
		}

		var versions []VersionMetaData
		if err := json.Unmarshal(fileData, &versions); err != nil {
			var singleVersion VersionMetaData
			if err := json.Unmarshal(fileData, &singleVersion); err != nil {
				return "", fmt.Errorf("failed to parse version data: %w", err)
			}
			versions = []VersionMetaData{singleVersion}
		}

		if len(versions) > 0 {
			highestVersion := 0
			for _, version := range versions {
				if len(version.ID) > 1 && version.ID[0] == 'v' {
					vNum, err := strconv.Atoi(version.ID[1:])
					if err == nil && vNum > highestVersion {
						highestVersion = vNum
					}
				}
			}
			nextVersion = highestVersion + 1
		}
	}

	versionID := fmt.Sprintf("v%d", nextVersion)
	return versionID, nil
}
