package version

import "time"

type VersionMetaData struct {
	ID        string
	Message   string
	Size      int64
	Checksum  string
	CreatedAt time.Time
}

type GlobalIndex struct {
	OriginalFilePath string
	Versions         []string
	LastUpdatedAt    time.Time
}
