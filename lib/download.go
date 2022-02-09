package lib

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

var skipVerify bool = false

func DownloadKctl(version string, path string, verify bool) (bool, error) {
	skipVerify = verify
	kctlFileLocation := fmt.Sprintf("%s/kubectl.%s", path, version)

	if err := checkPath(version, kctlFileLocation); err != nil {
		if err := downloadFile(version, kctlFileLocation); err != nil {
			libLogger.Errorf("Can't download kubectl version %s", version)
			return false, err
		}
	} else {
		libLogger.Infof("Found a binary for version %s with a valid checksum, skipping download.", version)
		return false, nil
	}
	return true, nil
}

func downloadFile(version string, path string) error {
	url := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", version, runtime.GOOS, runtime.GOARCH)
	libLogger.Info(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		libLogger.Errorf("Can't build a download client")
		return err
	}

	defer resp.Body.Close()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		libLogger.Errorf("Can't write destination file at %s.", f.Name())
		return err
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("downloading kubectl %s", version),
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	if err := verifyKctlDownload(version, path); err != nil {
		return err
	}
	return nil
}

func checkPath(version string, path string) error {

	dirPath := filepath.Dir(path)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) { // file does not exist
			if err := os.MkdirAll(dirPath, 0755); err != nil { // try to create path
				libLogger.Fatalf("error creating destination directory: " + err.Error()) // crash and burn if you can't
			}
		}
		return err
	} else { // file exists...
		if err := verifyKctlDownload(version, path); err != nil {
			libLogger.Debug(err.Error()) // ...but checksum don't match :(
			return err
		}
		return nil
	}
}

func verifyKctlDownload(version string, path string) error {
	if skipVerify {
		libLogger.Info("skipping file verification")
		return nil
	}
	checksumURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl.sha512", version, runtime.GOOS, runtime.GOARCH)
	libLogger.Info(checksumURL)

	resp, err := http.Get(checksumURL)
	if err != nil {
		libLogger.Errorf("Can't download sha512 checksums for version %s from %s.", version, checksumURL)
		return err
	}
	bodyData, _ := ioutil.ReadAll(resp.Body)
	checksumRef := strings.TrimSuffix(string(bodyData), "\n")

	f, err := os.Open(path)
	if err != nil {
		// TODO: don't Fatal
		libLogger.Fatal(err)
	}
	defer f.Close()

	h := sha512.New()
	if _, err := io.Copy(h, f); err != nil {
		return err
	}
	checksumCalc := hex.EncodeToString(h.Sum(nil))

	if checksumCalc != checksumRef {
		return fmt.Errorf("expected checksum %s, instead got %s", checksumRef, checksumCalc)
	}

	return nil
}
