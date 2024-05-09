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

// (WIP) Extracts the data from a .zip file to the specified target directory.
func unzip(reader io.ReadCloser, target string) error {
    buf := bytes.NewBuffer([]byte{})
    size, err := io.Copy(buf, reader)
    if err != nil {
        return err
    }

    r := bytes.NewReader(buf.Bytes())
    zipReader, err := zip.NewReader(r, size)
    if err != nil {
        return err
    }

    for _, file := range zipReader.File {
        destPath := filepath.Clean(filepath.Join(target, file.Name))
        if !strings.HasPrefix(destPath, target) {
            return errors.New("invalid file path")
        }
    }

    return nil
}

// Extracts the data from a .tar.gz file to the specified target directory. 
func untar(reader io.ReadCloser, target string) error {
    stream, err := gzip.NewReader(reader)
    if err != nil {
        return err
    }

    tarReader := tar.NewReader(stream)
    for {
        header, err := tarReader.Next()
        if err == io.EOF {
            break
        } else if err != nil {
            return err
        }

        path := filepath.Join(target, header.Name)
        info := header.FileInfo()
        if info.IsDir() {
            if err = os.MkdirAll(path, info.Mode()); err != nil {
                return err
            }

            continue
        }

        file, err := os.OpenFile(path, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, info.Mode())
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
