package schedule

import (
	"log"
	"time"
)

type QueryRequest struct {
	Query string
	Start time.Time
	End   time.Time
}

func Run() {
	queryConf := &QueryConf{
		Query: "{job=\"fluent-bit\",app=\"yinglong\"}",
		Schedule: Schedule{
			interval: intervalDefault,
			max:      maxDefault},
		Prefix:      "hello",
		NamePattern: "log-%s",
	}

	requests := generateRequests(queryConf)
	log.Printf("requests: %d", len(requests))
}

func generateRequests(conf *QueryConf) []QueryRequest {

	return nil
}
