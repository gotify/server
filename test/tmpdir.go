package test

import (
	"io/ioutil"
	"os"
	"path"
)

// TmpDir is a handler to temporary directory.
type TmpDir struct {
	path string
}

// Path returns the path to the temporary directory joined by the elements provided.
func (c TmpDir) Path(elem ...string) string {
	return path.Join(append([]string{c.path}, elem...)...)
}

// Clean removes the TmpDir.
func (c TmpDir) Clean() error {
	return os.RemoveAll(c.path)
}

// NewTmpDir returns a new handle to a tmp dir.
func NewTmpDir(prefix string) TmpDir {
	dir, _ := ioutil.TempDir("", prefix)
	return TmpDir{dir}
}
