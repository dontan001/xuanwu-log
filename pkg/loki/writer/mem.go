package writer

import (
	"log"
)

type MemWriter struct {
	Store []byte
}

func (w MemWriter) Write(p []byte) (n int, err error) {
	w.Store = append(w.Store, p...)

	i := 0
	for i < len(p) {
		log.Printf("%s", string(p[i]))
	}
	return len(p), nil
}
