package schedule

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
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
			S3: &s3.S3Config{
				Bucket: "donggetest",
				Region: endpoints.UsWest2RegionID,
			},
		},
	}

	store = func() *storage.Store {
		s := &storage.Store{
			Config: &s3.S3Config{
				Bucket: testConf.Archive.S3.Bucket,
				Region: testConf.Archive.S3.Region},
		}
		s.Setup()

		return s
	}()
)

func TestEnsure(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
	}
}

func TestGenerateRequests(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
		requests := query.generateRequests(store)

		log.Printf("total: %d", len(requests))
		for idx, request := range requests {
			log.Printf("Request #%d %s", idx+1, request.String())
		}
	}
}

func TestSubmit(t *testing.T) {
	for _, query := range testConf.Queries {
		query.ensure(testConf)
		requests := query.generateRequests(store)
		submit(requests)
	}
}
