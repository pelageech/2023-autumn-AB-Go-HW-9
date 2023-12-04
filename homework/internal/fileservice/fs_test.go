package fileservice_test

import (
	"bytes"
	"context"
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"homework/internal/fileservice"
	"homework/internal/models"
	"homework/pkg/iterator"
)

type mockFS = fstest.MapFS

var ctx = context.Background()
var testfs = mockFS{
	"bin/internal/usr/game": &fstest.MapFile{
		Data: []byte("super game!"),
		Mode: 0o766,
	},
	"bin/internal/ls.exe": &fstest.MapFile{
		Data: []byte("2"),
		Mode: 0o777,
	},
	"bin/usr": &fstest.MapFile{
		Mode: 0o777 | fs.ModeDir,
	},
	"bin/internal/file.txt": &fstest.MapFile{
		Data: []byte("6"),
		Mode: 0o147,
	},
}

func TestService_Ls(t *testing.T) {
	type args struct {
		path models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    []models.FileName
		wantErr bool
	}{
		{"file_not_exist_error", args{path: "bin/aboba"}, nil, true},
		{"no_files_OK", args{path: "bin/usr"}, []models.FileName{}, false},
		{"files_exist_OK", args{path: "bin/internal"}, []models.FileName{"usr", "ls.exe", "file.txt"}, false},
		{"file_ls_error", args{path: "bin/internal/ls.exe"}, nil, true},
	}
	s := fileservice.New(testfs)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Ls(ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Ls() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			for i := 0; i < len(got); i++ {
				got[i] = strings.TrimRight(got[i], "/")
			}
			if !assert.ElementsMatch(t, got, tt.want) {
				t.Errorf("Ls() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Meta(t *testing.T) {
	type args struct {
		path models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    *models.FileInfo
		wantErr bool
	}{
		{"file_not_exist_error", args{"bin/aboba"}, nil, true},
		{"file_OK", args{"bin/internal/usr/game"}, &models.FileInfo{
			Size:  11,
			Mode:  0o766,
			IsDir: false,
		}, false},
		{"folder_OK", args{"bin/internal"}, &models.FileInfo{
			Size:  0,
			IsDir: true,
		}, false},
	}

	s := fileservice.New(testfs)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Meta(ctx, tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Meta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil && (*tt.want != *got) {
				t.Errorf("Meta() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_ReadFile(t *testing.T) {
	type args struct {
		path models.FilePath
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"file_not_exist_error", args{"bin/aboba"}, nil, true},
		{"file_OK", args{"bin/internal/usr/game"}, []byte("super game!"), false},
		{"folder_error", args{"bin/internal"}, nil, true},
	}
	s := fileservice.New(testfs)

	buf := bytes.NewBuffer([]byte{})
	for _, tt := range tests {
		buf.Reset()
		t.Run(tt.name, func(t *testing.T) {
			i, err := s.ReadFileIterator(ctx, tt.args.path)
			if err == nil {
				err = iterator.Iterate(i, func(b []byte) error {
					_, _ = buf.Write(b)
					return nil
				})
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got := buf.Bytes()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
