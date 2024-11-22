//go:build !unix

package runner

func umask(_ int) int {
	return 0
}
