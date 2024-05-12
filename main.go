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

    tea "github.com/charmbracelet/bubbletea"
)

const (
    gopaDir     string = ".gopa"
    goDir       string = "go"
    logFile     string = "gopa.log"
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

    appDir = filepath.Join(homeDir, gopaDir)

    err = os.Mkdir(appDir, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        slog.Error("os.Mkdir", "error", err.Error())
        os.Exit(1)
    }

    file, err := os.Create(filepath.Join(appDir, logFile))
    if err != nil {
        slog.Error("os.Create", "error", err.Error())
        os.Exit(1)
    }
    defer file.Close()

    _, err = os.Stat(filepath.Join(appDir, goDir))
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
        filepath.Join(appDir, logFile), os.O_APPEND | os.O_WRONLY, os.ModePerm)
    if err != nil {
        slog.Error("os.OpenFile", "error", err.Error())
        os.Exit(1)
    }
    defer file.Close()

    logger := slog.New(
        slog.NewTextHandler(
            file,
            &slog.HandlerOptions{Level: logLevels[os.Getenv("LOG_LEVEL")]},
        ),
    )

    logger.Debug("Successfully started the application")

    if _, err := tea.NewProgram(newPlayground(), tea.WithAltScreen()).Run(); err != nil {
        logger.Error("tea.NewProgram", "error", err.Error())
        os.Exit(1)
    }
}
