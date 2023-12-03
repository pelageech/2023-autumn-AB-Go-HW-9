package server

import (
	"io/fs"

	"homework/internal/models"
)

type Service interface {
	ReadFile(path models.FilePath) ([]byte, error)
	Ls(path models.FilePath) ([]models.FileName, error)
	Meta(path models.FilePath) (fs.FileInfo, error)
}
