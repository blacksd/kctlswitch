package lib_test

import (
	"errors"
	"fmt"
	"io/fs"
	"kctlswitch/lib"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ValidKubectlVersion string = "1.12.3"

var validSrcPath string

var InstallTests = []struct {
	name    string
	version string
	srcPath string
	dstPath string
	want    error
}{
	{"valid dstPath without binaries", ValidKubectlVersion, validSrcPath, "/usr/local/bin/", nil},
	{"valid dstPath with symlink", ValidKubectlVersion, validSrcPath, "./", nil},
	{"valid dstPath with non-symlink binary", ValidKubectlVersion, validSrcPath, "/usr/bin/", errors.New("won't override an existing binary")},
	{"empty dstPath", ValidKubectlVersion, validSrcPath, "", &fs.PathError{Op: "stat", Path: "", Err: syscall.Errno(2)}},
	{"dstPath is not a dir", ValidKubectlVersion, validSrcPath, "/bin/bash", errors.New("destination path is not a directory")},
	{"non-existing dstPath", ValidKubectlVersion, validSrcPath, "/probablythisdoesnot/exists/", &fs.PathError{Op: "stat", Path: "/probablythisdoesnot/exists/", Err: syscall.Errno(2)}},
	{"inaccessible dstpath", ValidKubectlVersion, validSrcPath, "/bin/", &os.LinkError{Op: "symlink", Old: fmt.Sprintf("%s.v%s", lib.DefaultKctlBinaryName, ValidKubectlVersion), New: "/bin/kubectl", Err: syscall.Errno(1)}},
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Fatal("Can't get the current user home directory. Stopping here.")
	}
	validSrcPath = fmt.Sprintf("%s/.kctlswitch/bin/kubectl-v%s", userHomeDir, ValidKubectlVersion)
	lib.DownloadKctl(ValidKubectlVersion, validSrcPath, slog)
	//TODO: implement test PASS dependency on KctlDownload
	//TODO: initialize tests
	/*
		- call KctlDownload on ValidKubectlVersion to validSrcPath
		- build local testdata
	*/
}

func TestInstall(t *testing.T) {
	for _, it := range InstallTests {
		t.Run(it.name, func(t *testing.T) {
			err := lib.InstallKctlVersion(it.version, it.srcPath, it.dstPath, slog)
			assert.Equal(t, it.want, err)
			if it.want == nil {
				l, err := os.Readlink(it.srcPath)
				if assert.Nil(t, err) {
					assert.Equal(t, l, it.dstPath)
				}
			}
		})
	}
}
