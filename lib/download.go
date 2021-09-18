package lib

import (
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/schollz/progressbar/v3"
)

func DownloadKctl(kctlVersion string, installLocation string) error {

	// TODO: check sha256
	// TODO: check if file is already there
	if path, err := validatePath(installLocation); err == nil {
		kctlFileLocation := fmt.Sprintf("%s/kubectl", path)
		// if err := downloadFile(kctlVersion, kctlFileLocation); err != nil {
		// 	log.Fatalf("Can't download kubectl version %s", kctlVersion)
		// }
		if err := verifyKctlDownload(kctlVersion, kctlFileLocation); err != nil {
			log.Fatal("Checksum computation error.")
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

func verifyKctlDownload(kctlVersion string, kctlFileLocation string) error {
	checksumURL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl.sha512", kctlVersion, runtime.GOOS, runtime.GOARCH)
	log.Print(checksumURL)

	resp, err := http.Get(checksumURL)
	if err != nil {
		log.Printf("Can't download sha512 checksums for version %s from %s.", kctlVersion, checksumURL)
	}
	bodyData, _ := ioutil.ReadAll(resp.Body)
	checksumRef := strings.TrimSuffix(string(bodyData), "\n")

	f, err := os.Open(kctlFileLocation)
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
