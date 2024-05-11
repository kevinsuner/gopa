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

Objectives (11 May 2024)
- [x] Application initialization (log file and golang installation)
- [x] Save debug, info, warn and error logs to log file depending on the LOG_LEVEL
- [ ] Make the terminal display a textbox where the user can input text
    - [ ] Integration test: Initialize the program, write to the textbox, check the textbox input and quit the program
- [ ] Make a panel to the right-side of the textbox, that echoes the content of the textbox on <CTRL-Enter>
    - [ ] Integration test: Initialize the program, write to the textbox, press <CTRL-Enter>, match textbox data against
panel data, quit the program

*/

const (
    APP_DIR_NAME    string = ".gopa"
    GO_DIR_NAME     string = "go"
    LOG_FILE_NAME   string = "gopa.log"
)

var (
    // Logging levels mapped to their slog.Level representation.
    logLevels = map[string]slog.Level{
        "":         slog.LevelInfo, // default logging level
        "debug":    slog.LevelDebug,
        "info":     slog.LevelInfo,
        "warn":     slog.LevelWarn,
        "error":    slog.LevelError,
    }

    // Format string to download latest Golang version based on OS and ARCH.
    goURL = "https://go.dev/dl/go1.22.3.%s-%s.%s"

    // Path including the user's home to the .gopa directory
    appDir string
)

func init() {
    homeDir, err := os.UserHomeDir()
    if err != nil {
        slog.Error("os.UserHomeDir", "error", err.Error())
        os.Exit(1)
    }

    appDir = filepath.Join(homeDir, APP_DIR_NAME)

    err = os.Mkdir(appDir, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        slog.Error("os.Mkdir", "error", err.Error())
        os.Exit(1)
    }

    file, err := os.Create(filepath.Join(appDir, LOG_FILE_NAME))
    if err != nil {
        slog.Error("os.Create", "error", err.Error())
        os.Exit(1)
    }
    defer file.Close()

    _, err = os.Stat(filepath.Join(appDir, GO_DIR_NAME))
    if os.IsNotExist(err) {
        ext := "tar.gz"
        if runtime.GOOS == "windows" { ext = "zip" }

        resp, err := http.Get(fmt.Sprintf(goURL, runtime.GOOS, runtime.GOARCH, ext))
        if err != nil {
            slog.Error("http.Get", "error", err.Error())
            os.Exit(1)
        }

        if err := Uncompress(resp.Body, appDir, runtime.GOOS); err != nil {
            slog.Error("Uncompress", "error", err.Error())
            os.Exit(1)
        }
    }
}

func main() {
    file, err := os.OpenFile(
        filepath.Join(appDir, LOG_FILE_NAME), os.O_APPEND | os.O_WRONLY, os.ModePerm)
    if err != nil {
        panic(err)
    }
    defer file.Close()

    logger := slog.New(
        slog.NewTextHandler(
            file,
            &slog.HandlerOptions{Level: logLevels[os.Getenv("LOG_LEVEL")]},
        ),
    )

    logger.Debug("Successfully started application")
}
