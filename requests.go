/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/mod/semver"
)

const reqURL string = "https://go.dev"
var ErrUnexpectedStatus error = errors.New("unexpected status code")

func getLatestGoVersion() (string, error) {
    res, err := http.Get(fmt.Sprintf("%s/VERSION?m=text", reqURL))
    if err != nil {
        return "", err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return "", ErrUnexpectedStatus 
    }

    body, err := io.ReadAll(res.Body)
    if err != nil {
        return "", err
    }

    return strings.Split(string(body), "\n")[0], nil
}

func downloadGoVersion(dst, version, osys string) error {
    ext := "tar.gz"
    if osys == "windows" { ext = "zip" }

    _, err := os.Stat(filepath.Join(dst, version))
    if os.IsNotExist(err) {
        res, err := http.Get(fmt.Sprintf("%s/dl/%s.%s", reqURL, version, ext))
        if err != nil { return err }
        defer res.Body.Close()

        if res.StatusCode != http.StatusOK {
            return ErrUnexpectedStatus
        }

        if err := uncompress(res.Body, dst, runtime.GOOS); err != nil {
            return err
        }

        return os.Rename(filepath.Join(dst, "go"), filepath.Join(dst, version))
    }

    return nil
}

func listGoVersions() ([]string, error) {
    res, err := http.Get(fmt.Sprintf("%s/dl", reqURL))
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode != http.StatusOK {
        return nil, ErrUnexpectedStatus
    }

    doc, err := goquery.NewDocumentFromReader(res.Body)
    if err != nil {
        return nil, err
    }

    rawVersions := []string{}
    doc.Find(".toggleButton").Each(func(i int, s *goquery.Selection) {
        version := s.Find("span").Text()

        // Match versions from go1.16+ ahead and replace the leading "go"
		// prefix for a "v" prefix to sort it using the semver package
		regex := regexp.MustCompile(`^go(\d+)\.(1[6-9]|[2-9]\d+)(?:\.(\d+))?$`)
		if regex.MatchString(version) {
			rawVersions = append(rawVersions, strings.Replace(version, "go", "v", 1))
		}
    })

    rawVersions = slices.Compact(rawVersions)
    semver.Sort(rawVersions)
    slices.Reverse(rawVersions)

    versions := []string{}
    for _, rawVersion := range rawVersions {
        versions = append(versions, strings.Replace(rawVersion, "v", "go", 1))
    }

    return versions, nil
}

