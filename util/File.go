package util

import (
	"os"
	"syscall"
)

func _PathExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Hide(path string) error {
	if !_PathExists(path) {
		return nil
	}
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return syscall.SetFileAttributes(cpath, syscall.FILE_ATTRIBUTE_HIDDEN)
}
