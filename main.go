/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"log/slog"
	"os"
	"path/filepath"
)

/*

Primary objectives
- [x] Application initialization (log file and golang installation)
- [x] Save debug, info, warn and error logs to log file depending on the LOG_LEVEL
    - [ ] Integration test: Initialize the desired directories and create logs with different LOG_LEVEL's
- [ ] Make the terminal display a textbox where the user can input text
    - [ ] Integration test: Initialize the program, write to the textbox, check the textbox input and quit the program

*/


// Human-readable logging levels mapped to their slog.Level representation.
var logLevels = map[string]slog.Level{
    "":         slog.LevelInfo, // default logging level
    "debug":    slog.LevelDebug,
    "info":     slog.LevelInfo,
    "warn":     slog.LevelWarn,
    "error":    slog.LevelError,
}

func main() {
    rootDir, err := Init()
    if err != nil {
        panic(err)
    }

    file, err := os.OpenFile(
        filepath.Join(rootDir, LOG_FILE), os.O_APPEND | os.O_WRONLY, os.ModePerm)
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
