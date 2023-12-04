package iterator

import (
	"context"
	"errors"
	"io"
)

// ReaderIterator is an iterator that reads io.Reader by chunks.
// It implements Simple.
type ReaderIterator struct {
	r   io.Reader
	eof bool
	buf []byte
	ctx context.Context
}

// NewReaderIterator created new ReaderIterator.
func NewReaderIterator(ctx context.Context, r io.Reader, bufSize int) *ReaderIterator {
	return &ReaderIterator{
		r:   r,
		buf: make([]byte, bufSize),
		ctx: ctx,
	}
}

func (i *ReaderIterator) HasNext() bool {
	return !i.eof
}

func (i *ReaderIterator) Next() ([]byte, error) {
	n, err := i.r.Read(i.buf)
	if errors.Is(err, io.EOF) {
		i.eof = true
	} else if err != nil {
		return nil, err
	}

	return i.buf[:n], i.ctx.Err()
}

func (i *ReaderIterator) Reader() io.Reader {
	return i.r
}
