package query

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func TestQueryStdout(t *testing.T) {
	result := os.Stdout
	Query("{job=\"fluent-bit\"}", time.Now().Add(-10*time.Minute), time.Now(), result)
	t.Logf("done")
}

func TestQueryBuffer(t *testing.T) {
	// result := &writer.MemWriter{}
	result := &bytes.Buffer{}
	Query("{job=\"fluent-bit\"}", time.Now().Add(-10*time.Minute), time.Now(), result)
	t.Logf("result: %d", len(result.Bytes()))
	t.Logf("result: %s", result.String())
}
