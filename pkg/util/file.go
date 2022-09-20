package util

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func Concatenate(merged, other string) error {
	create := false
	if _, err := os.Stat(merged); errors.Is(err, os.ErrNotExist) {
		create = true
	}

	target, err := os.OpenFile(merged, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer target.Close()

	src, err := os.Open(other)
	if err != nil {
		return err
	}
	defer src.Close()

	if !create {
		target.WriteString("\n")
	}

	_, err = io.Copy(target, src)
	return nil
}

func Compress(source, target string) error {
	return zipSource(source, target)
}

func zipSource(source, target string) error {
	defer TimeMeasure("zip")()

	// Create a ZIP file and zip.Writer
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := zip.NewWriter(f)
	defer writer.Close()

	// Go through all the files of the source
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Create a local file header
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		// set compression
		header.Method = zip.Deflate

		// Set relative path of a file as the header name
		header.Name, err = filepath.Rel(filepath.Dir(source), path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			header.Name += "/"
		}

		// Create writer for the file header and save content of the file
		headerWriter, err := writer.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(headerWriter, f)
		return err
	})
}

func UnCompress(source, dest string) ([]string, error) {
	fileNames, err := unzip(source, dest)
	return fileNames, err
}

func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {
		fPath := filepath.Join(dest, f.Name)

		prefix := filepath.Clean(dest) + string(os.PathSeparator)
		if !strings.HasPrefix(fPath, prefix) {
			return filenames, fmt.Errorf("%s is an illegal filepath", fPath)
		}

		filenames = append(filenames, fPath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return filenames, err
		}

		err = copy(fPath, f)
		if err != nil {
			return filenames, err
		}
	}

	return filenames, nil
}

func copy(fPath string, f *zip.File) error {
	outFile, err := os.OpenFile(fPath,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()

	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	_, err = io.Copy(outFile, rc)
	if err != nil {
		return err
	}

	return nil
}

func RandomName() string {

	id, _ := uuid.NewUUID()
	return id.String()
}
