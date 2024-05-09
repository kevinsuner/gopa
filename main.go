/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
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

// Start creates the required folders and files for the project to work,
// and downloads the latest Go version based on the host's operating system.
func Start() error {
    home, err := os.UserHomeDir()
    if err != nil {
        return err
    }

    rootDir := filepath.Join(home, APP_DIR)
    err = os.Mkdir(rootDir, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        return err
    }

    file, err := os.Create(filepath.Join(rootDir, LOG_FILE))
    if err != nil {
        return err
    }

    defer func() error {
        return file.Close()
    }()

    _, err = os.Stat(filepath.Join(rootDir, GO_DIR))
    if os.IsNotExist(err) {
        ext := "tar.gz"
        if runtime.GOOS == "windows" { ext = "zip" }

        resp, err := http.Get(fmt.Sprintf(goURL, runtime.GOOS, runtime.GOARCH, ext))
        if err != nil {
            return err
        }

        if err := Uncompress(resp.Body, rootDir, runtime.GOOS); err != nil {
            return err
        }
    }

    return nil
}

func main() {
    log := slog.New(
        slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

    if err := Start(); err != nil {
        log.Error("Start()", "error", err)
    }

    log.Debug("Successfully started application")
}
