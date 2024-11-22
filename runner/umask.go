//go:build unix

package runner

import "syscall"

var umask = syscall.Umask
