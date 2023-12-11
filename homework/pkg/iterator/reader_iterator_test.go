package iterator

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"
)

func TestReaderIterator(t *testing.T) {
	type fields struct {
		r       io.Reader
		bufSize int
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{"empty reader", fields{bytes.NewReader(nil), 1}, []byte{}, false},
		{"one iter", fields{bytes.NewReader([]byte{1, 2, 3, 4, 5}), 10}, []byte{1, 2, 3, 4, 5}, false},
		{"multi iter", fields{bytes.NewReader([]byte{1, 2, 3, 4, 5}), 1}, []byte{1, 2, 3, 4, 5}, false},
	}

	buf := bytes.NewBuffer([]byte{})
	for _, tt := range tests {
		bufIterator := make([]byte, tt.fields.bufSize)

		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			i := NewReaderIterator(context.Background(), tt.fields.r, bufIterator)
			err := Iterate[[]byte](i, func(b []byte) error {
				buf.Write(b)
				return nil
			})

			if (err != nil) != tt.wantErr {
				t.Errorf("Next() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := buf.Bytes()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Next() got = %v, want %v", got, tt.want)
			}
		})
	}
}
