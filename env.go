/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"fmt"
	"os"
)

func setGoBin(path, osys string) error {
    exe := "go"
    if osys == "windows" { exe = "go.exe" }

    return os.Setenv("GOPA_GO_BIN", fmt.Sprintf("%s/bin/%s", path, exe))
}

func getGoBin() string {
    return os.Getenv("GOPA_GO_BIN")
}
