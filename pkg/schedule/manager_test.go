package schedule

import (
	"github.com/kyligence/xuanwu-log/pkg/data/loki"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
)

var (
	testBackup = &Backup{
		Data: &data.DataConf{
			Loki: &loki.LokiConf{
				Address: "http://aafdd592dddec49ed8bf3c35d9d538c9-577636166.us-west-2.elb.amazonaws.com:80",
			},
		},
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

	testData = func() *data.Data {
		d := &data.Data{Conf: testBackup.Data}
		d.Setup()

		return d
	}()

	testStore = func() *storage.Store {
		s := &storage.Store{
			Config: &s3.S3Config{
				Bucket: testBackup.Archive.S3.Bucket,
				Region: testBackup.Archive.S3.Region},
		}
		s.Setup()

		return s
	}()
)

func TestEnsure(t *testing.T) {
	for _, query := range testBackup.Queries {
		query.Ensure(BACKUP, testBackup)
	}
}

func TestGenerateRequests(t *testing.T) {
	for _, query := range testBackup.Queries {
		query.Ensure(BACKUP, testBackup)
		requests := query.generateRequests(testData, testStore)

		log.Printf("total: %d", len(requests))
		for idx, request := range requests {
			log.Printf("Request #%d %s", idx+1, request.String())
		}
	}
}

func TestSubmit(t *testing.T) {
	for _, query := range testBackup.Queries {
		query.Ensure(BACKUP, testBackup)
		requests := query.generateRequests(testData, testStore)
		submit(requests)
	}
}
