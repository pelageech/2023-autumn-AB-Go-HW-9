package fileservice_test

import (
	"io/fs"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/stretchr/testify/assert"

	"homework/internal/fileservice"
	"homework/internal/models"
)

type mapFileInfo struct {
	name string
	f    *fstest.MapFile
}

func (i *mapFileInfo) Name() string               { return i.name }
func (i *mapFileInfo) Size() int64                { return int64(len(i.f.Data)) }
func (i *mapFileInfo) Mode() fs.FileMode          { return i.f.Mode }
func (i *mapFileInfo) Type() fs.FileMode          { return i.f.Mode.Type() }
func (i *mapFileInfo) ModTime() time.Time         { return i.f.ModTime }
func (i *mapFileInfo) IsDir() bool                { return i.f.Mode&fs.ModeDir != 0 }
func (i *mapFileInfo) Sys() any                   { return i.f.Sys }
func (i *mapFileInfo) Info() (fs.FileInfo, error) { return i, nil }

func (i *mapFileInfo) String() string {
	return fs.FormatFileInfo(i)
}

type mockFS = fstest.MapFS

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
			got, err := s.Ls(tt.args.path)
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
		want    fs.FileInfo
		wantErr bool
	}{
		{"file_not_exist_error", args{"bin/aboba"}, nil, true},
		{"file_OK", args{"bin/internal/usr/game"}, &mapFileInfo{
			name: "game",
			f: &fstest.MapFile{
				Data: []byte("super game!"),
				Mode: 0o766,
			},
		}, false},
		{"folder_OK", args{"bin/internal"}, &mapFileInfo{
			name: "internal",
			f: &fstest.MapFile{
				Mode: 0o000 | fs.ModeDir,
			},
		}, false},
	}

	s := fileservice.New(testfs)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.Meta(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Meta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr != true && !(got.Name() == tt.want.Name() &&
				got.IsDir() == tt.want.IsDir() &&
				got.Size() == tt.want.Size() &&
				got.Mode() == tt.want.Mode()) {
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.ReadFile(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}
