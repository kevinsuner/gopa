/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
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
    gosDirname   string = "gos"
    logFilename string = "gopa.log"
)

var (
    rootDir string
    goURL   string = "https://go.dev"
)

func init() {
    home, err := os.UserHomeDir()
    if err != nil {
        slog.Error("os.UserHomeDir", "error", err.Error())
        os.Exit(1)
    }

    rootDir = filepath.Join(home, appDirname)
    err = os.Mkdir(rootDir, os.ModePerm)
    if err != nil && !errors.Is(err, fs.ErrExist) {
        slog.Error("os.Mkdir", "error", err.Error())
        os.Exit(1)
    }

    file, err := os.Create(filepath.Join(rootDir, logFilename))
    if err != nil {
        slog.Error("os.Create", "error", err.Error())
        os.Exit(1)
    }
    defer file.Close()

    version, err := getLatestGoVersion()
    if err != nil {
        slog.Error("getLatestGoVersion", "error", err.Error())
        os.Exit(1)
    }

    _, err = os.Stat(filepath.Join(rootDir, gosDirname))
    if os.IsNotExist(err) {
        err = os.Mkdir(filepath.Join(rootDir, gosDirname), os.ModePerm)
        if err != nil {
            slog.Error("os.Stat", "error", err.Error())
            os.Exit(1)
        }
    }

    longVersion := fmt.Sprintf("%s.%s-%s", version, runtime.GOOS, runtime.GOARCH)
    _, err = os.Stat(filepath.Join(rootDir, gosDirname, longVersion))
    if os.IsNotExist(err) {
        ext := "tar.gz"
        if runtime.GOOS == "windows" { ext = "zip" }

        if err := downloadGoVersion(longVersion, ext); err != nil {
            slog.Error("downloadGoVersion", "error", err.Error())
            os.Exit(1)
        }
    }

    if err := os.Setenv("GOPA_GO_VERSION", longVersion); err != nil {
        slog.Error("os.Setenv", "error", err.Error())
        os.Exit(1)
    }
}

func getLatestGoVersion() (string, error) {
    resp, err := http.Get(fmt.Sprintf("%s/VERSION?m=text", goURL))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", errors.New("unexpected status code")
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return strings.Split(string(body), "\n")[0], nil
}

func downloadGoVersion(version, ext string) error {
    resp, err := http.Get(
        fmt.Sprintf("%s/dl/%s.%s", goURL, version, ext))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    err = Uncompress(
        resp.Body, filepath.Join(rootDir, gosDirname), runtime.GOOS)
    if err != nil {
        return err
    }

    err = os.Rename(
        filepath.Join(rootDir, gosDirname, "go"),
        filepath.Join(rootDir, gosDirname, version))
    if err != nil {
        return err
    }

    return nil
}

type playground struct {
    editor *tview.TextArea
    console *tview.TextView
    menu *tview.Box
    flex *tview.Flex
}

func newPlayground() playground {
    console := newConsole()
    editor := newEditor(console)
    menu := newMenu()

    return playground{
        editor: editor,
        console: console,
        menu: menu,
        flex: tview.NewFlex().
            AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
                AddItem(tview.NewFlex().SetDirection(tview.FlexColumn).
                    AddItem(editor, 0, 1, false).
                    AddItem(console, 0, 1, false), 0, 1, false).
                AddItem(menu, 5, 1, false), 0, 1, false),
    }
}

func newConsole() *tview.TextView {
    console := tview.NewTextView().SetWordWrap(true)
    console.SetTitle("Console").SetBorder(true)
    return console
}

func newEditor(console *tview.TextView) *tview.TextArea {
    editor := tview.NewTextArea().
        SetWrap(false).
        SetPlaceholder("Type some code here...")
    
    editor.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
        if event.Key() == tcell.KeyCtrlR {
            out, err := runCode(editor.GetText())
            if err != nil {
                panic(err)
            }
            
            console.SetText(out)
            return nil
        }

        return event
    })

    editor.SetTitle("Editor").SetBorder(true)
    return editor
}

func runCode(input string) (string, error) {
    file, err := createTempFile()
    if err != nil {
        return "", err
    }

    defer file.Close()
    defer os.Remove(file.Name())

    _, err = file.WriteString(input)
    if err != nil {
        return "", err
    }

    goExec := filepath.Join(os.Getenv("GOPA_GO_VERSION"), "bin", "go")
    if runtime.GOOS == "windows" {
        goExec = filepath.Join(os.Getenv("GOPA_GO_VERSION"), "bin", "go.exe") }

    cmd := exec.Command(filepath.Join(rootDir, gosDirname, goExec), "run", file.Name())
    cmd.Dir = rootDir  
    out, _ := cmd.CombinedOutput()

    return string(out), nil
}

func createTempFile() (*os.File, error) {
    return os.CreateTemp(rootDir, "gopa-*.go")
}

func newMenu() *tview.Box {
    return tview.NewBox().SetBorder(true)
}

func main() {
    app := tview.NewApplication()
    playground := newPlayground()
    if err := app.SetRoot(playground.flex, true).EnableMouse(true).EnablePaste(true).Run(); err != nil {
        panic(err)
    }
}
