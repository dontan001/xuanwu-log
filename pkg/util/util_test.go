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
	raw := "%7Bapp=%22yl-common-booter%22,namespace=%22yinglong-dev%22,job=%22fluent-bit%22%7D"
	unescaped, err := url.QueryUnescape(raw)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", unescaped)

	raw = "{app=\"yl-common-booter\",namespace=\"yinglong-dev\",job=\"fluent-bit\"}"
	unescaped, err = url.QueryUnescape(raw)
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%v", unescaped)
}
