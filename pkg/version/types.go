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

type DiffResult struct {
	Identical bool
	DiffType  string
	Message   string
	DiffLines []LineDiff
}

type LineDiff struct {
	LineNumber int
	Line1      string
	Line2      string
}
