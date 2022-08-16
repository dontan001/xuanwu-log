package util

import (
	"testing"
	"time"
)

func TestParseTimeExpression(t *testing.T) {
	timeParsed, e := ParseTime("now+1h")
	if e != nil {
		t.Fatalf("%s", e)
	}
	t.Logf("%s", timeParsed.Format(time.RFC3339))
}

func TestParseTimeNow(t *testing.T) {
	timeParsed, e := ParseTime("now")
	if e != nil {
		t.Fatalf("%s", e)
	}
	t.Logf("%s", timeParsed.Format(time.RFC3339))
}

func TestParseTimeFormat(t *testing.T) {
	_, e := ParseTime("now11")
	if e != nil {
		t.Logf("%s", e)
	}
}
