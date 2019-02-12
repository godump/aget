package aget

import (
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type OpenFunc func(string) (io.ReadCloser, error)

func OpenFile(name string) (io.ReadCloser, error) {
	return os.Open(name)
}

func OpenHTTP(name string) (io.ReadCloser, error) {
	resp, err := http.Get(name)
	if err != nil {
		return nil, err
	}
	return resp.Body, nil
}

func Open(name string) (io.ReadCloser, error) {
	switch {
	case strings.HasPrefix(name, "http://"), strings.HasPrefix(name, "https://"):
		return OpenHTTP(name)
	}
	return OpenFile(name)
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
		return os.Open(save)
	}
	return f(name)
}

func OpenFileEx(name string, save string, ex time.Duration) (io.ReadCloser, error) {
	return withEx(OpenFile, name, save, ex)
}

func OpenHTTPEx(name string, save string, ex time.Duration) (io.ReadCloser, error) {
	return withEx(OpenHTTP, name, save, ex)
}

func OpenEx(name string, save string, ex time.Duration) (io.ReadCloser, error) {
	return withEx(Open, name, save, ex)
}
