package util

import (
	"testing"
	"time"
)

func TestZipSource(t *testing.T) {
	fs := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/pkg/util/zip.go"
	ft := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/pkg/util/zip.go.zip"
	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
}

func TestZipSourceBig(t *testing.T) {
	start := time.Now()
	fs := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/test50g.txt"
	ft := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/test50g.txt.zip"
	if err := ZipSource(fs, ft); err != nil {
		t.Fatal(err)
	}
	t.Logf("elapsed %f seconds", time.Since(start).Seconds())
}
