package util

import (
	"os"
	"syscall"
)

type IFile struct {
	path       string
	fileObject os.File
}

func _IFileExists(ifile IFile) bool {
	path := ifile.path
	return _PathExists(path)
}
func _PathExists(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func Hide(ifile IFile) error {
	if !_IFileExists(ifile) {
		return nil
	}
	path := ifile.path
	cpath, cpathErr := syscall.UTF16PtrFromString(path)
	if cpathErr != nil {
		return cpathErr
	}
	return syscall.SetFileAttributes(cpath, syscall.FILE_ATTRIBUTE_HIDDEN)
}
