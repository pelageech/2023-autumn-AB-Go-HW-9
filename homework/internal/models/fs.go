package models

type (
	FileName = string
	FilePath = string
)

type FileInfo struct {
	Size  int64
	Mode  uint32
	IsDir bool
}
