/*
* SPDX-License-Identifier: GPL-3.0-only
* Copyright (C) 2024 Kevin Suñer <ksuner@pm.me>
 */

package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Uncompress reads a stream of data, and determines how to uncompress it based
// on the operating system passed in the function arguments.
func Uncompress(reader io.ReadCloser, target, os string) error {
    if os == "windows" {
        return unzip(reader, target)
    }

    return untar(reader, target)
}

// Extracts the data from a .zip file to the specified target directory.
func unzip(reader io.ReadCloser, target string) error {
    buf := bytes.NewBuffer([]byte{})
    size, err := io.Copy(buf, reader)
    if err != nil {
        return err
    }

    bytesReader := bytes.NewReader(buf.Bytes())
    zipReader, err := zip.NewReader(bytesReader, size)
    if err != nil {
        return err
    }

    for _, file := range zipReader.File {
        path := filepath.Clean(filepath.Join(target, file.Name))
        if !strings.HasPrefix(path, target) {
            return errors.New("invalid file path")
        }
    
        if file.FileInfo().IsDir() {
            if err := os.MkdirAll(path, os.ModePerm); err != nil { return err }
            continue
        }

        if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
            return err
        }

        dstFile, err := os.OpenFile(
            path, os.O_CREATE | os.O_TRUNC | os.O_RDWR, file.Mode())
        if err != nil {
            return err
        }
        defer dstFile.Close()

        srcFile, err := file.Open()
        if err != nil {
            return err
        }
        defer srcFile.Close()

        if _, err := io.Copy(dstFile, srcFile); err != nil {
            return err
        }
    }

    return nil
}

// Extracts the data from a .tar.gz file to the specified target directory.
func untar(reader io.ReadCloser, target string) error {
    gzipReader, err := gzip.NewReader(reader)
    if err != nil {
        return err
    }
    
    defer func() error {
        return gzipReader.Close()
    }()

    tarReader := tar.NewReader(gzipReader)
    for {
        header, err := tarReader.Next()
        if err != nil {
            if err == io.EOF { break }
            return err
        }

        path := filepath.Join(target, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err := os.MkdirAll(path, info.Mode()); err != nil { return err }
            continue
        }

        file, err := os.OpenFile(
            path, os.O_CREATE | os.O_TRUNC | os.O_RDWR, info.Mode())
        if err != nil {
            return err
        }
        defer file.Close()

        _, err = io.Copy(file, tarReader)
        if err != nil {
            return err
        }
    }

    return nil
}

