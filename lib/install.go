package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"golang.org/x/sys/unix"
)

const (
	defaultKctlBinaryName        string = "kubectl"
	defaultKctlBinaryVersionsDir string = ".kube/bin/"
)

var (
	ErrNotADir               error = errors.New("path is not a directory")
	ErrNotSymlinkFilePresent error = fmt.Errorf("there's already a %s in path, and it's not a symlink; please set the --force flag to overwrite it", defaultKctlBinaryName)
)

func InstallKctlVersion(kctlVersion string, srcBinaryPath string, dstSymlinkPath string, overwrite bool) error {
	if err := validateDstSymlinkPath(dstSymlinkPath, overwrite); err != nil {
		libLogger.Error(err)
		return err
	}

	srcBinary := filepath.Join(srcBinaryPath, fmt.Sprintf("%s.v%s", defaultKctlBinaryName, kctlVersion))
	dstSymlink := filepath.Join(dstSymlinkPath, defaultKctlBinaryName)

	if currentDestination, err := os.Lstat(dstSymlink); err == nil {
		libLogger.Debug("found a symlink pointing to %s", currentDestination.Name())
		if err := os.Remove(dstSymlink); err != nil {
			libLogger.Error("failed to unlink the existing link")
			return err
		}
	}

	if err := os.Symlink(srcBinary, dstSymlink); err != nil {
		return err
	}
	libLogger.Info("symlink successfully set")
	return nil
}

func validateDstSymlinkPath(path string, overwrite bool) error {
	// check the path is valid (name and it's a directory)
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return ErrNotADir
	}

	// path is not writable
	if err := unix.Access(path, unix.W_OK); err != nil {
		return err
	}

	defaultBinary := filepath.Join(path, defaultKctlBinaryName)
	if binInfo, err := os.Stat(defaultBinary); err == nil {
		if (binInfo.Mode()&os.ModeSymlink != os.ModeSymlink) && !overwrite {
			return ErrNotSymlinkFilePresent
		}
	} else {
		if err.(*os.PathError).Err.(syscall.Errno) != syscall.ENOENT {
			return err
		}
	}

	return nil
}
