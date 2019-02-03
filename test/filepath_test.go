package test

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProjectPath(t *testing.T) {
	_, err := os.Stat(path.Join(GetProjectDir(), "./README.md"))
	assert.Nil(t, err)
}

func TestWithWd(t *testing.T) {
	wd1, _ := os.Getwd()
	tmpDir := NewTmpDir("gotify_withwd")
	defer tmpDir.Clean()
	var wd2 string
	WithWd(tmpDir.Path(), func(origWd string) {
		assert.Equal(t, wd1, origWd)
		wd2, _ = os.Getwd()
	})
	wd3, _ := os.Getwd()
	assert.Equal(t, wd1, wd3)
	assert.Equal(t, tmpDir.Path(), wd2)
	assert.Nil(t, os.RemoveAll(tmpDir.Path()))

	assert.Panics(t, func() {
		WithWd("non_exist", func(string) {})
	})

	assert.Nil(t, os.Mkdir(tmpDir.Path(), 0644))
	assert.Panics(t, func() {
		WithWd(tmpDir.Path(), func(string) {})
	})
	assert.Nil(t, os.Remove(tmpDir.Path()))

	assert.Nil(t, os.Mkdir(tmpDir.Path(), 0755))
	assert.Panics(t, func() {
		WithWd(tmpDir.Path(), func(string) {
			assert.Nil(t, os.RemoveAll(tmpDir.Path()))
			WithWd(".", func(string) {})
		})
	})

}
