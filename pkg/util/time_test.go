package util

import (
	"testing"
	"time"
)

func TestParseTimeUnixNano(t *testing.T) {
	timeParsed, e := ParseTime("1661203728487614000")
	if e != nil {
		t.Fatalf("%s", e)
	}
	t.Logf("%d", timeParsed.UnixNano())
	t.Logf("%s", timeParsed.Format(time.RFC3339Nano))
}

func TestParseTimeExpression(t *testing.T) {
	timeParsed, e := ParseTime("now+1h")
	if e != nil {
		t.Fatalf("%s", e)
	}
	t.Logf("%s", timeParsed.Format(time.RFC3339Nano))
}

func TestParseTimeNow(t *testing.T) {
	timeParsed, e := ParseTime("now")
	if e != nil {
		t.Fatalf("%s", e)
	}
	t.Logf("%s", timeParsed.Format(time.RFC3339Nano))
}

func TestParseTimeFormat(t *testing.T) {
	_, e := ParseTime("now11")
	if e != nil {
		t.Logf("%s", e)
	}
}

func TestCalcLastBackup(t *testing.T) {
	interval := 3
	tm := time.Now()

	lbt := CalcLastBackup(interval, tm)
	t.Logf("%s", lbt.Format(time.RFC3339Nano))
}
