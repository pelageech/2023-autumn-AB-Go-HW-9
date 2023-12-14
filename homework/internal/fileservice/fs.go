package fileservice

import (
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

// RepositoryFS is an equivalent to the file system interface.
type RepositoryFS interface {
	fs.ReadDirFS
}

// Service contains methods for FS calls.
type Service struct {
	fs RepositoryFS
}

// New creates an implementation of a new Service.
func New(fs RepositoryFS) *Service {
	return &Service{fs: fs}
}

// ReadFileIterator opens a file on FS and returns an iterator *iterator.ReaderIterator.
// Note that the method doesn't close the file.
//
// The iterator contains the reader that implements io.ReadCloser.
// The user should close it manually.
func (s *Service) ReadFileIterator(ctx context.Context, path models.FilePath) (_ *iterator.ReaderIterator, err error) {
	f, err := s.fs.Open(path)
	if err != nil {
		return nil, fmt.Errorf("file open error: %w", err)
	}

	return iterator.NewReaderIterator(ctx, f, make([]byte, packetSize)), nil
}

// Ls returns a list of files containing in the given path.
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

// Meta returns meta-data of the file in path if exists.
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
