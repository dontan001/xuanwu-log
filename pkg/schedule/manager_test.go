package schedule

import (
	"log"
	"testing"
)

var (
	testConf = &BackupConf{
		Queries: []*QueryConf{
			{
				Query: "{job=\"fluent-bit\",app=\"yinglong\"}",
				Schedule: &Schedule{
					Interval: DefaultInterval,
					Max:      DefaultMax},
			},
			{
				Query: "{job=\"fluent-bit\",app=\"yinglong\",node=\"ip-10-1-254-253.us-west-2.compute.internal\"}",
				Schedule: &Schedule{
					Interval: DefaultInterval,
					Max:      DefaultMax},
			},
		},
		Archive: &Archive{
			Type:        DefaultType,
			WorkingDir:  DefaultWorkingDir,
			NamePattern: "%s.log",
		},
	}
)

func TestEnsure(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
	}
}

func TestGenerateRequests(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
		requests := query.generateRequests()

		log.Printf("total: %d", len(requests))
		for idx, request := range requests {
			log.Printf("Request #%d %s", idx+1, request.String())
		}
	}
}

func TestSubmit(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
		requests := query.generateRequests()
		submit(requests)
	}
}
