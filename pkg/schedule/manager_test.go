package schedule

import (
	"log"
	"testing"
)

var (
	testConf = &QueryConf{
		Query: "{job=\"fluent-bit\",app=\"yinglong\"}",
		Schedule: Schedule{
			Interval: DefaultInterval,
			Max:      DefaultMax},
		Prefix:      "test",
		NamePattern: "test-%s",
	}
)

func TestGenerateRequests(t *testing.T) {
	requests := generateRequests(testConf)
	log.Printf("total: %d", len(requests))
	for idx, request := range requests {
		log.Printf("Request #%d %s", idx+1, request.String())
	}
}

func TestSubmit(t *testing.T) {
	requests := generateRequests(testConf)
	submit(requests)
}
