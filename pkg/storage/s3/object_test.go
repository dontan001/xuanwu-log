package s3

import (
	"fmt"
	"testing"
)

const (
	BASE = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/%s"
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
	err := HeadObject(remotePath)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetObject(t *testing.T) {
	remotePath := "index/loki_index_19240/loki-loki-distributed-ingester-0-1662344787333739953-1662348480.gz"
	fileName := fmt.Sprintf(BASE, "test/tmp.txt")
	err := GetObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestGetObjectBig(t *testing.T) {
	remotePath := "test/test1g.txt"
	fileName := fmt.Sprintf(BASE, "test/tmp.txt")
	err := GetObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestPutObject(t *testing.T) {
	fileName := fmt.Sprintf(BASE, "test/README.md")
	remotePath := "test/README.md"

	err := PutObject(remotePath, fileName)
	if err != nil {
		t.Logf("%s", err)
	}
}

func TestPutObjectBig(t *testing.T) {
	fileName := fmt.Sprintf(BASE, "test/test1g.txt")
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
