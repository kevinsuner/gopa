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

/*

Primary objectives
- [ ] Application initialization (log file and golang installation)
- [ ] Save debug, info, warn and error logs to log file depending on the LOG_LEVEL
    - [ ] Integration test: Initialize the desired directories and create logs with different LOG_LEVEL's
- [ ] Make the terminal display a textbox where the user can input text
    - [ ] Integration test: Initialize the program, write to the textbox, check the textbox input and quit the program

*/

const (
    APP_DIR     string = ".gopa"
    LOG_FILE    string = "gopa.log"
    GO_DIR      string = "go"
)

var goURL string = "https://go.dev/dl/go1.22.3.%s-%s.%s"

// Initializes the required folder/files required for the project to work.
func init() {
    home, err := os.UserHomeDir()
    if err != nil {
        // TODO: log don't panic
        panic(err)
    }

    err = os.Mkdir(filepath.Join(home, APP_DIR), os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        // TODO: log don't panic
        panic(err)
    }

    file, err := os.Create(filepath.Join(home, APP_DIR, LOG_FILE))
    if err != nil {
        // TODO: log don't panic
        panic(err)
    }
    defer file.Close()

    _, err = os.Stat(filepath.Join(home, APP_DIR, GO_DIR))
    if os.IsNotExist(err) {
        ext := "tar.gz"
        if runtime.GOOS == "windows" { ext = "zip" }

        resp, err := http.Get(fmt.Sprintf(goURL, runtime.GOOS, runtime.GOARCH, ext))
        if err != nil {
            // TODO: log don't panic
            panic(err)
        }

        if runtime.GOOS == "windows" {
            // TODO: Implement unzip functionallity for windows
            return
        }

        err = untar(resp.Body, filepath.Join(home, APP_DIR))
        if err != nil {
            // TODO: log don't panic
            panic(err)
        }
    }
}

func main() {
    fmt.Println("Hello from Gopa!")
}
