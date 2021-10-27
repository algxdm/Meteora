package util

import (
	"os"
	"os/exec"
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

func HideWin(path string) error {
	cmd := exec.Command("attrib", "+H", path, "/S", "/L")
	return cmd.Run()
}
