package query

import (
	"bytes"
	"os"
	"testing"
	"time"
)

const (
	queryStr = "{job=\"fluent-bit\"}"
)

func TestQueryStdout(t *testing.T) {
	result := os.Stdout
	Query(queryStr, time.Now().Add(-10*time.Minute), time.Now(), result)
	t.Logf("done")
}

func TestQueryBuffer(t *testing.T) {
	result := &bytes.Buffer{}
	Query(queryStr, time.Now().Add(-10*time.Minute), time.Now(), result)
	t.Logf("result: %d", len(result.Bytes()))
	t.Logf("result: %s", result.String())
}
