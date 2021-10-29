package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"github.com/spf13/afero"
	"go.uber.org/zap"
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

func InstallKctlVersion(kctlVersion string, srcPath string, dstPath string, overwrite bool, AFS afero.Fs, log *zap.SugaredLogger) error {
	if err := validateDstPath(dstPath, overwrite, AFS); err != nil {
		log.Error(err)
		return err
	}

	srcBinary := filepath.Join(srcPath, fmt.Sprintf("%s.v%s", defaultKctlBinaryName, kctlVersion))
	dstBinary := filepath.Join(dstPath, defaultKctlBinaryName)

	// if currentDestination, err := os.Lstat(dstBinary); err == nil {

	// if currentDestination, _, err := AFS.(*afero.OsFs).LstatIfPossible(dstBinary); err == nil {
	if currentDestination, _, err := AFS.(afero.Lstater).LstatIfPossible(dstBinary); err == nil {
		log.Debug("found a symlink pointing to %s", currentDestination.Name())
		if err := AFS.Remove(dstBinary); err != nil {
			log.Error("failed to unlink the existing link")
			return err
		}
	}

	if err := AFS.(afero.Symlinker).SymlinkIfPossible(srcBinary, dstBinary); err != nil {
		return err
	}
	log.Info("symlink successfully set")
	return nil
}

func validateDstPath(path string, overwrite bool, AFS afero.Fs) error {
	// check the path is valid (name and it's a directory)
	pathInfo, err := AFS.Stat(path)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return ErrNotADir
	}

	// path is not writable
	// BUG: fix this, it's not going to work with afero
	if err := unix.Access(path, unix.W_OK); err != nil {
		return err
	}

	defaultBinary := filepath.Join(path, defaultKctlBinaryName)
	if binInfo, err := AFS.Stat(defaultBinary); err == nil {
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
