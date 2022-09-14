package s3

import (
	"path/filepath"
	"testing"
)

const (
	WorkingDir = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test"
)

func TestGetBuckets(t *testing.T) {
	err := GetBuckets()
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetObjects(t *testing.T) {
	err := GetObjects()
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestHeadObject(t *testing.T) {
	remotePath := "index/loki_index_19240/loki-loki-distributed-ingester-0-1662344787333739953-1662348480.gz"
	_, err := HeadObject(remotePath)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestHeadObject404(t *testing.T) {
	remotePath := "index/loki_index_19240/loki-loki-distributed-ingester-0-1662344787333739953-xxxx.gz"
	_, err := HeadObject(remotePath)
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetObject(t *testing.T) {
	remotePath := "index/loki_index_19240/loki-loki-distributed-ingester-0-1662344787333739953-1662348480.gz"
	fileName := filepath.Join(WorkingDir, "tmp.txt")
	err := GetObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetObjectBig(t *testing.T) {
	remotePath := "test/test1g.txt"
	fileName := filepath.Join(WorkingDir, "tmp.txt")
	err := GetObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestPutObject(t *testing.T) {
	fileName := filepath.Join(WorkingDir, "README.md")
	remotePath := "test/README.md"

	err := PutObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestPutObjectBig(t *testing.T) {
	fileName := filepath.Join(WorkingDir, "test1g.txt")
	remotePath := "test/test1g.txt"

	err := PutObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestDelObject(t *testing.T) {
	remotePath := "test/client.go"

	err := DelObject(remotePath)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestDelObject404(t *testing.T) {
	remotePath := "test/xxx.go"

	err := DelObject(remotePath)
	if err != nil {
		t.Logf("%s", err)
	}
}
