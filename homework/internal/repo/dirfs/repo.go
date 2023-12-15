package dirfs

import (
	"io/fs"
	"os"
)

func New(dir string) fs.ReadDirFS {
	return os.DirFS(dir).(fs.ReadDirFS) // the result implements fs.ReadDirFS
}
