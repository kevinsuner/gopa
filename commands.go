/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"os"
	"os/exec"
)

func runCode(input string) (string, error) {
    file, err := os.CreateTemp(appDir, "gopa-*.go")
    if err != nil {
        return "", err
    }
    defer file.Close()
    defer os.Remove(file.Name())

    _, err = file.WriteString(input)
    if err != nil {
        return "", err
    }

    cmd := exec.Command(getGoBin(), "run", file.Name())
    cmd.Dir = appDir
    out, _ := cmd.CombinedOutput()

    return string(out), nil
}

