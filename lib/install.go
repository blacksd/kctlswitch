package lib

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
)

const (
	defaultKctlBinaryName        string = "kubectl"
	defaultKctlBinaryVersionsDir string = ".kube/bin/"
)

var (
	ErrNotADir error = errors.New("path is not a directory")
)

func InstallKctlVersion(kctlVersion string, srcPath string, dstPath string, log *zap.SugaredLogger) error {

	if err := validateDstPath(dstPath); err != nil {
		log.Error(err)
		return err
	}

	srcBinary := filepath.Join(srcPath, fmt.Sprintf("%s.v%s", defaultKctlBinaryName, kctlVersion))
	dstBinary := filepath.Join(dstPath, defaultKctlBinaryName)

	// TODO: this check should happen before, in the validateDstPath
	// TODO: set an explicit --force flag to overwrite destination, or unset link instead
	if currentDestination, err := os.Lstat(dstBinary); err == nil {
		log.Debug("found a symlink pointing to %s", currentDestination.Name())
		if err := os.Remove(dstBinary); err != nil {
			log.Error("failed to unlink the existing link")
			return err
		}
	}

	if err := os.Symlink(srcBinary, dstBinary); err != nil {
		return err
	}
	log.Info("symlink successfully set")
	return nil
}

func validateDstPath(path string) error {
	// check the path is valid (name and it's a directory)
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !pathInfo.IsDir() {
		return ErrNotADir
	}

	// check we can write there
	// check there's no other non-symlink file

	// f, err := os.OpenFile(filepath.Join(path, defaultKctlBinaryName), os.O_CREATE, 0644)
	// if err != nil {
	// 	return err
	// 	//os.IsPermission(err)
	// }
	// f.Close()

	// TODO: test for O_CREATE|O_WRITE perms
	// TODO: check if there's already a kubectl, and if it is a symlink

	return nil
}
