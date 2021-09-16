package lib

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

func DownloadKctl(kctlVersion string, installLocation string) error {

	URL := fmt.Sprintf("https://storage.googleapis.com/kubernetes-release/release/%s/bin/%s/%s/kubectl", kctlVersion, runtime.GOOS, runtime.GOARCH)
	log.Print(URL)
	// TODO: check sha256
	if err := downloadFile(URL, fmt.Sprintf("%s/kubectl", installLocation)); err != nil {
		log.Fatalf("Can't download file at %s", URL)
	}

	return nil
}

func downloadFile(url string, path string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// func verifyKctlDownload(installLocation string) error {
//
// 	return nil
// }
