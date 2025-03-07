package version

func CreateFile(filePath, versionID, message string) (VersionMetaData, error) {
	meta, err := saveFile(filePath, versionID, message)
	return meta, err
}
