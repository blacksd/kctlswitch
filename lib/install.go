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

func InstallKctlVersion(kctlVersion string, srcPath string, dstPath string, log *zap.SugaredLogger) error {

	if err := validateDstPath(dstPath); err != nil {
		log.Error(err)
		return err
	}

	srcBinary := filepath.Join(srcPath, fmt.Sprintf("%s.v%s", defaultKctlBinaryName, kctlVersion))
	dstBinary := filepath.Join(dstPath, defaultKctlBinaryName)

	if _, err := os.Lstat(dstBinary); err == nil {
		log.Debug("found a symlink pointing to XX")
		if err := os.Remove(dstBinary); err != nil {
			log.Error("failed to unlink the existing link")
			return err
		}
	}

	if err := os.Symlink(srcBinary, dstBinary); err != nil {
		return err
	}
	log.Info("Symlink successfully set")
	return nil
}

func validateDstPath(path string) error {
	// TODO: check the path is valid and we can write there
	pathInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if !pathInfo.IsDir() {
		// TODO: refactor this error to be a var
		return errors.New("destination path is not a directory")
	}

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
