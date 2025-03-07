package version

import (
	"errors"
)

func CreateFile(filePath, versionID, message string) (VersionMetaData, error) {
	fileDir, err := getVersionPath(filePath)
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
