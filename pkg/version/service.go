package version

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func CreateFile(filePath, versionID, message string) (VersionMetaData, error) {
	fileDir, err := GetVersionPath(filePath)
	if err != nil {
		return VersionMetaData{}, err
	}
	lastFilePath := returnLastFilePath(fileDir)
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
func ListAllVersions(fileDir string) ([]VersionMetaData, error) {
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
	return listAllVersions, nil
}

///////////////////////////////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////////////////////////////
