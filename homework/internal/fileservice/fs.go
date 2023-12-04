package fileservice

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io/fs"

	"homework/internal/models"
	"homework/pkg/iterator"
)

// packetSize is 4KiB - a size which write() operation
// can process in one step (writing a system block).
const packetSize = 4 << 10 // 4 KiB

type ReadDirFS interface {
	fs.ReadDirFS
}

type Service struct {
	fs ReadDirFS
}

func New(fs ReadDirFS) *Service {
	return &Service{fs: fs}
}

func (s *Service) ReadFileIterator(ctx context.Context, path models.FilePath) (_ iterator.Interface[[]byte], err error) {
	f, err := s.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	buf := bufio.NewReader(f)

	return iterator.Reader(ctx, buf, packetSize), nil
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

func (s *Service) Meta(ctx context.Context, path models.FilePath) (_ *models.FileInfo, err error) {
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

	info := &models.FileInfo{
		Size:  stat.Size(),
		Mode:  uint32(stat.Mode().Perm()),
		IsDir: stat.IsDir(),
	}

	return info, nil
}
