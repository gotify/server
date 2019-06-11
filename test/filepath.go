package test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
)

// GetProjectDir returns the correct absolute path of this project
func GetProjectDir() string {
	_, f, _, _ := runtime.Caller(0)
	projectDir, _ := filepath.Abs(path.Join(filepath.Dir(f), "../"))
	return projectDir
}

// WithWd executes a function with the specified working directory
func WithWd(chDir string, f func(origWd string)) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if err := os.Chdir(chDir); err != nil {
		panic(err)
	}
	defer os.Chdir(wd)
	f(wd)
}
