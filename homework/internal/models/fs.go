package models

type (
	FileName = string
	FilePath = string
)

// FileInfo contains information about files
type FileInfo struct {
	Size  int64
	Mode  uint32
	IsDir bool
}
