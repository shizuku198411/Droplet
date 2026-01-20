package container

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// mount path validation for protecting path traversal
func securePath(rootfs, dest string) (string, error) {
	rel := strings.TrimPrefix(dest, "/")
	clean := filepath.Clean(rel)

	if clean == "." || strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || clean == ".." {
		return "", fmt.Errorf("invalid destination: %q", dest)
	}

	fullPath := filepath.Join(rootfs, clean)
	root := filepath.Clean(rootfs) + string(os.PathSeparator)
	if !strings.HasPrefix(fullPath+string(os.PathSeparator), root) && fullPath != filepath.Clean(rootfs) {
		return "", fmt.Errorf("destination escapes rootfs: %q -> %q", dest, fullPath)
	}
	return fullPath, nil
}

// source mount validation
// the following source is denied by default
//
//	/proc, /sys, /dev, /run, /var/run, /boot, /root, /
func hasDeniedSource(source string) bool {
	p := filepath.Clean(source)
	if p == "/" {
		return true
	}
	deniedList := []string{"/proc", "/sys", "/dev", "/run", "/var/run", "/boot", "/root"}
	for _, d := range deniedList {
		if p == d || strings.HasPrefix(p, d+string(os.PathSeparator)) {
			return true
		}
	}
	return false
}

func isSymlink(source string) (bool, error) {
	fi, err := os.Lstat(source)
	if err != nil {
		return false, err
	}
	return fi.Mode()&os.ModeSymlink != 0, nil
}

func secureMount(source, target, fstype string, flags uintptr, data string) error {
	// 1. bind mount
	if err := syscall.Mount(source, target, fstype, flags, data); err != nil {
		return err
	}

	// 2. change mount type
	if err := syscall.Mount("", target, "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return err
	}
	return nil
}
