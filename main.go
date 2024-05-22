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
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Ideas
// - Think of switching UIs as an array of UIs ui[0]

const (
    appDirname  string = ".gopa"
    gosDirname  string = "gos"
    logFilename string = "gopa.log"
    gosFilename string = ".gos"
)

var (
    appDir      string
    goVersions  []string
)

func setup() error {
    home, err := os.UserHomeDir()
    if err != nil { return err }

    appDir = filepath.Join(home, appDirname)
    err = os.Mkdir(appDir, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        return err
    }

    file, err := os.Create(filepath.Join(appDir, logFilename))
    if err != nil { return err }
    defer file.Close()

    _, err = os.Stat(filepath.Join(appDir, gosDirname))
    if os.IsNotExist(err) {
        err := os.Mkdir(filepath.Join(appDir, gosDirname), os.ModePerm)
        if err != nil { return err }
    }

    version, err := getLatestGoVersion()
    if err != nil { return err }

    longVersion := fmt.Sprintf("%s.%s-%s", version, runtime.GOOS, runtime.GOARCH)
    _, err = os.Stat(filepath.Join(appDir, gosDirname, longVersion))
    if os.IsNotExist(err) {
        err = downloadGoVersion(
            filepath.Join(appDir, gosDirname), longVersion, runtime.GOOS)   
        if err != nil {
            return err
        }
    }

    goVersions, err = cacheGoVersions()
    if err != nil { return err }

    err = setGoBin(filepath.Join(appDir, gosDirname, longVersion), runtime.GOOS)
    if err != nil { return err }

    return nil
}

func main() {
    if err := setup(); err != nil {
        slog.Error("setup", "error", err.Error())
        os.Exit(1)
    }

    playground := newPlayground(tview.NewApplication())
    playground.layout = playground.newLayout()
    playground.editor = playground.extendEditor()
    playground.console = playground.extendConsole()
    playground.menu = playground.extendMenu()
    playground.versionsList = playground.newVersionsList()

    playground.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlL {
            for _, version := range goVersions {
                playground.versionsList.AddItem(
                    version,
                    fmt.Sprintf("Go version %s", strings.Replace(version, "go", "", 1)),
                    '\n',
                    func() {
                        longVersion := fmt.Sprintf("%s.%s-%s", version, runtime.GOOS, runtime.GOARCH)
                        err := downloadGoVersion(
                            filepath.Join(appDir, gosDirname), longVersion, runtime.GOOS)   
                        if err != nil {
                            slog.Error("downloadGoVersion", "error", err.Error())
                            os.Exit(1)
                        }

                        err = setGoBin(filepath.Join(appDir, gosDirname, longVersion), runtime.GOOS)
                        if err != nil {
                            slog.Error("setGoBin", "error", err.Error())
                            os.Exit(1)
                        }

                        playground.app.SetRoot(playground.layout, true)
                    },
                )
            }

            playground.app.SetRoot(playground.versionsList, true)
            return nil
        }

        return event
    })

    err := playground.app.SetRoot(playground.layout, true).
        EnableMouse(true).
        EnablePaste(true).
        Run()
    if err != nil {
        slog.Error("playground.app.Run", "error", err.Error())
        os.Exit(1)
    }
}

