package iterator

import (
	"context"
	"errors"
	"io"
)

type ReaderIterator struct {
	r   io.Reader
	eof bool
	buf []byte
	ctx context.Context
}

func Reader(ctx context.Context, r io.Reader, bufSize int) *ReaderIterator {
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
