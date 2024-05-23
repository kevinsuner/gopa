/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func cacheGoVersions() ([]string, error) {
    _, err := os.Stat(filepath.Join(appDir, gosFilename))
    if os.IsNotExist(err) {
        file, err := os.Create(filepath.Join(appDir, gosFilename))
        if err != nil {
            return nil, err
        }
        defer file.Close()

        versions, err := listGoVersions()
        if err != nil {
            return nil, err
        }

        _, err = file.WriteString(
            fmt.Sprintf("%s\n", time.Now().Format(time.DateOnly)))
        if err != nil {
            return nil, err
        }

        for _, version := range versions {
            _, err := file.WriteString(fmt.Sprintf("%s\n", version))
            if err != nil {
                return nil, err
            }
        }

        return versions, nil
    }

    file, err := os.OpenFile(
        filepath.Join(appDir, gosFilename), os.O_RDWR, os.ModePerm)
    if err != nil {
        return nil, err
    }

    scanner := bufio.NewScanner(file)
    scanner.Split(bufio.ScanLines)
    
    var idx int
    var expired bool
    versions := []string{}

    for scanner.Scan() {
        if idx == 0 {
            timestamp, err := time.Parse(time.DateOnly, scanner.Text())
            if err != nil {
                return nil, err
            }

            now, err := time.Parse(time.DateOnly, time.Now().Format(time.DateOnly))
            if err != nil {
                return nil, err
            }

            if timestamp.Before(now) { expired = true; break }
            idx++; continue
        }

        versions = append(versions, scanner.Text())
    }

    if expired {
        versions, err = listGoVersions()
        if err != nil {
            return nil, err
        }

        file.Truncate(0)
        file.Seek(0, 0)

        _, err = file.WriteString(
            fmt.Sprintf("%s\n", time.Now().Format(time.DateOnly)))
        if err != nil {
            return nil, err
        }

        for _, version := range versions {
            _, err := file.WriteString(fmt.Sprintf("%s\n", version))
            if err != nil {
                return nil, err
            }
        }
    }

    return versions, nil
}
