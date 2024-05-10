/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

const (
    APP_DIR     string = ".gopa"
    LOG_FILE    string = "gopa.log"
    GO_DIR      string = "go"
)

// Format string to download latest Golang version based on OS and ARCH.
var goURL = "https://go.dev/dl/go1.22.3.%s-%s.%s"

// Init creates the required folders and files for the project to work,
// and downloads the latest Go version based on the host's operating system.
func Init(path string) (string, error) {
    basepath := filepath.Join(path, APP_DIR)

    err := os.Mkdir(basepath, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        return "", err
    }

    file, err := os.Create(filepath.Join(basepath, LOG_FILE))
    if err != nil {
        return "", err
    }
    defer file.Close()

    _, err = os.Stat(filepath.Join(basepath, GO_DIR))
    if os.IsNotExist(err) {
        ext := "tar.gz"
        if runtime.GOOS == "windows" { ext = "zip" }

        resp, err := http.Get(fmt.Sprintf(goURL, runtime.GOOS, runtime.GOARCH, ext))
        if err != nil {
            return "", err
        }

        if err := Uncompress(resp.Body, basepath, runtime.GOOS); err != nil {
            return "", err
        }
    }

    return basepath, nil
}

