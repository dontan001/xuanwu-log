package util

import (
	"path/filepath"
	"testing"
)

const (
	WorkingDir = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log"
)

func TestZipSource(t *testing.T) {
	fs := filepath.Join(WorkingDir, "pkg/util/file.go")
	ft := filepath.Join(WorkingDir, "pkg/util/file.go.zip")

	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
}

func TestZipSourceBig(t *testing.T) {
	defer TimeMeasure("zip")()

	fs := filepath.Join(WorkingDir, "test/test50g.txt")
	ft := filepath.Join(WorkingDir, "test/test50g.txt.zip")

	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
}

func TestConcatenate(t *testing.T) {
	ft := filepath.Join(WorkingDir, "test/t1.txt")
	fs := filepath.Join(WorkingDir, "test/t2.txt")

	if err := Concatenate(ft, fs); err != nil {
		t.Fatal(err)
	}
}

func TestUnzip(t *testing.T) {
	src := filepath.Join(WorkingDir, "test/util.go.zip")
	dest := filepath.Join(WorkingDir, "test")

	fileNames, err := Unzip(src, dest)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("unzip done: %v", fileNames)
}
