package util

import (
	"fmt"
	"testing"
)

const (
	BASE = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/%s"
)

func TestZipSource(t *testing.T) {
	fs := fmt.Sprintf(BASE, "pkg/util/zip.go")
	ft := fmt.Sprintf(BASE, "pkg/util/zip.go.zip")

	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
}

func TestZipSourceBig(t *testing.T) {
	defer TimeMeasure("zip")()

	fs := fmt.Sprintf(BASE, "test/test50g.txt")
	ft := fmt.Sprintf(BASE, "test/test50g.txt.zip")

	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
}
