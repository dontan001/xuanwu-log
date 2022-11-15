package util

import (
	"net/url"
	"testing"
)

func TestHash(t *testing.T) {
	t.Logf("%d", Hash("test"))
	t.Logf("%d", Hash("test."))
}

func TestDivMod(t *testing.T) {
	d, m := DivMod(15, 6)
	t.Logf("%d %d", d, m)

	d, m = DivMod(7, 6)
	t.Logf("%d %d", d, m)
}

func TestQueryUnescape(t *testing.T) {
	queryStr := "{job=\"fluent-bit\",app=\"yinglong\"}"
	escaped := url.QueryEscape(queryStr)
	t.Logf("%v", escaped)

	unescaped, err := url.QueryUnescape(escaped)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", unescaped)

	unescaped, err = url.QueryUnescape(unescaped)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", unescaped)
}
