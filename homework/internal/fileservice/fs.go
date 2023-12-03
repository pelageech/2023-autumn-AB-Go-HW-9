package fileservice

import (
	"bufio"
	"context"
	"errors"
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

func (s *Service) ReadFile(ctx context.Context, path models.FilePath) (_ []byte, err error) {
	defer func() {
		if ctx.Err() != nil {
			err = errors.Join(err, ctx.Err())
		}
	}()

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

func (s *Service) Ls(ctx context.Context, path models.FilePath) (_ []models.FileName, err error) {
	defer func() {
		if ctx.Err() != nil {
			err = errors.Join(err, ctx.Err())
		}
	}()

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

func (s *Service) Meta(ctx context.Context, path models.FilePath) (_ fs.FileInfo, err error) {
	defer func() {
		if ctx.Err() != nil {
			err = errors.Join(err, ctx.Err())
		}
	}()

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
