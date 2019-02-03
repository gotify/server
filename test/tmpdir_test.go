package test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTmpDir(t *testing.T) {
	dir := NewTmpDir("test_prefix")
	assert.NotEmpty(t, dir)

	assert.Contains(t, dir.Path(), "test_prefix")
	testFilePath := dir.Path("testfile.txt")
	assert.Contains(t, testFilePath, "test_prefix")
	assert.Contains(t, testFilePath, "testfile.txt")
	assert.True(t, strings.HasPrefix(testFilePath, dir.Path()))
}
