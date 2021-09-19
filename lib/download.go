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
	"go.uber.org/zap"
)

func DownloadKctl(version string, path string, log *zap.SugaredLogger) error {

	kctlFileLocation := fmt.Sprintf("%s/kubectl.%s", path, version)

	if err := checkPath(version, kctlFileLocation, log); err != nil {
		if err := downloadFile(version, kctlFileLocation, log); err != nil {
			log.Fatalf("Can't download kubectl version %s", version)
		}
	} else {
		log.Infof("Found a binary for version %s with the right checksum, skipping download.", version)
	}
	return nil
}

func downloadFile(version string, path string, log *zap.SugaredLogger) error {
	url := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", version, runtime.GOOS, runtime.GOARCH)
	log.Info(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Can't build a download client")
		return err
	}

	defer resp.Body.Close()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't write destination file at %s.", f.Name())
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("downloading kubectl %s", version),
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	if err := verifyKctlDownload(version, path, log); err != nil {
		return err
	}
	return nil
}

func checkPath(version string, path string, log *zap.SugaredLogger) error {

	dirPath := filepath.Dir(path)

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) { // file does not exist
			if err := os.MkdirAll(dirPath, 0755); err != nil { // try to create path
				log.Fatalf("error creating destination directory: " + err.Error()) // crash and burn if you can't
			}
		}
		return err
	} else { // file exists...
		if err := verifyKctlDownload(version, path, log); err != nil {
			log.Debug(err.Error()) // ...but checksum don't match :(
			return err
		}
		return nil
	}
}

func verifyKctlDownload(version string, path string, log *zap.SugaredLogger) error {
	checksumURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl.sha512", version, runtime.GOOS, runtime.GOARCH)
	log.Info(checksumURL)

	resp, err := http.Get(checksumURL)
	if err != nil {
		log.Errorf("Can't download sha512 checksums for version %s from %s.", version, checksumURL)
		return err
	}
	bodyData, _ := ioutil.ReadAll(resp.Body)
	checksumRef := strings.TrimSuffix(string(bodyData), "\n")

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
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
