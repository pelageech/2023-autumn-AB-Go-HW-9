package fileservice

import (
	"bufio"
	"fmt"
	"io"
	"io/fs"

	"homework/internal/models"
)

type ReadDirFS interface {
	fs.ReadDirFS
}

type Service struct {
	fs ReadDirFS
}

func New(fs ReadDirFS) *Service {
	return &Service{fs: fs}
}

func (s *Service) ReadFile(path models.FilePath) ([]byte, error) {
	f, err := s.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	buf := bufio.NewReader(f)
	b, err := io.ReadAll(buf)
	if err != nil {
		return nil, fmt.Errorf("file read error: %w", err)
	}

	return b, nil
}

func (s *Service) Ls(path models.FilePath) ([]models.FileName, error) {
	dir, err := s.fs.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir error: %w", err)
	}

	filenames := make([]models.FileName, len(dir))
	for k, v := range dir {
		filenames[k] = v.Name()
	}

	return filenames, nil
}

func (s *Service) Meta(path models.FilePath) (fs.FileInfo, error) {
	f, err := s.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	stat, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("file stat read error: %w", err)
	}

	return stat, nil
}
