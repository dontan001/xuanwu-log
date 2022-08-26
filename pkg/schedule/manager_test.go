package schedule

import (
	"log"
	"testing"
)

var (
	queryConf = &QueryConf{
		Query: "{job=\"fluent-bit\",app=\"yinglong\"}",
		Schedule: Schedule{
			Interval: DefaultInterval,
			Max:      DefaultMax},
		Prefix:      "test",
		NamePattern: "test-%s",
	}
)

func TestGenerateRequests(t *testing.T) {
	requests := generateRequests(queryConf)
	log.Printf("total: %d", len(requests))
	for idx, request := range requests {
		log.Printf("Request #%d %s", idx+1, request.String())
	}
}
