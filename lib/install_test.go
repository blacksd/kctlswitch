package lib_test

import (
	"errors"
	"fmt"
	"os"
)

const ValidKubectlVersion string = "v1.12.3"

var validSrcPath string

var InstallTests = []struct {
	name    string
	version string
	srcPath string
	dstPath string
	want    error
}{
	{"valid dstPath", ValidKubectlVersion, validSrcPath, "/usr/local/bin/", nil},
	{"empty dstPath", ValidKubectlVersion, validSrcPath, "./", errors.New("destination path can't be empty")},
	{"non-existing dstPath", ValidKubectlVersion, validSrcPath, "./", errors.New("destination path does not exist")},
	{"inaccessible dstpath", ValidKubectlVersion, validSrcPath, "./", errors.New("we can't write to destination path")},
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		slog.Fatal("Can't get the current user home directory. Stopping here.")
	}
	validSrcPath = fmt.Sprintf("%s/.kctlswitch/bin/kubectl-%s", userHomeDir, ValidKubectlVersion)
}

/* func TestInstall (t *testing.T) {
	for _, install := range validInstallTests {

		srcPath string =

	}
} */
