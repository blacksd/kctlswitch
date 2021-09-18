package lib

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/schollz/progressbar/v3"
)

func DownloadKctl(kctlVersion string, installLocation string) error {

	// TODO: check sha256
	// TODO: check if file is already there
	if path, err := validatePath(installLocation); err == nil {
		if err := downloadFile(kctlVersion, fmt.Sprintf("%s/kubectl", path)); err != nil {
			log.Fatalf("Can't download kubectl version %s", kctlVersion)
		}
	}

	return nil
}

func downloadFile(kctlVersion string, path string) error {
	url := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", kctlVersion, runtime.GOOS, runtime.GOARCH)
	log.Print(url)
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Print("Can't build a download client")
		return err
	}

	defer resp.Body.Close()

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Can't write destination file at %s", f.Name())
	}
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("downloading kubectl %s", kctlVersion),
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

	return nil
}

func validatePath(path string) (string, error) {
	if path == "" {
		path = "."
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path, 0755); err != nil {
				return "", errors.New("error creating destination directory: " + err.Error())
			}
		} else {
			return "", errors.New("error checking destination directory: " + err.Error())
		}
	}
	return path, nil
}

// func verifyKctlDownload(installLocation string) error {
//
// 	return nil
// }
