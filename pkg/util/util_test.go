package util

import (
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
