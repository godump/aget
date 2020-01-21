package aget

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// The OpenFunc type is an adapter to allow get an io.ReadCloser from
// specified path.
type OpenFunc func(string) (io.ReadCloser, error)

// OpenFile is an implemention of OpenFunc. It reads data from the disk.
func OpenFile(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

// OpenHTTP is an implemention of OpenFunc. It reads data by http.Get.
func OpenHTTP(name string) (io.ReadCloser, error) {
	resp, err := http.Get(name)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func isHTTP(name string) bool {
	return strings.HasPrefix(name, "http://") || strings.HasPrefix(name, "https://")
}

// Open is an implemention of OpenFunc. It select the appropriate method to
// open the file based on the incoming args automatically.
//
// Examples:
//   aget.Open("/etc/hosts")
//   aget.Open("https://github.com/godump/aget/blob/master/README.md")
func Open(name string) (io.ReadCloser, error) {
	switch {
	case isHTTP(name):
		return OpenHTTP(name)
	default:
		return OpenFile(name)
	}
}

func withEx(f OpenFunc, name string, save string, ex time.Duration) (io.ReadCloser, error) {
	i, err := os.Stat(save)
	if err != nil || time.Since(i.ModTime()) > ex {
		if err := func() error {
			wc, err := os.OpenFile(save, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				return err
			}
			defer wc.Close()
			rc, err := f(name)
			if err != nil {
				return err
			}
			defer rc.Close()
			_, err = io.Copy(wc, rc)
			return err
		}(); err != nil {
			return nil, err
		}
	}
	return os.Open(save)
}

// OpenHTTPEx is same as OpenHTTP but with cache.
func OpenHTTPEx(name string, save string, ex time.Duration) (io.ReadCloser, error) {
	return withEx(OpenHTTP, name, save, ex)
}

// OpenEx is same as Open but with cache.
func OpenEx(name string, save string, ex time.Duration) (io.ReadCloser, error) {
	switch {
	case isHTTP(name):
		return OpenHTTPEx(name, save, ex)
	default:
		return OpenFile(name)
	}
}
