package api

import (
	"log"
	"path/filepath"
	"testing"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/data/loki"
	"github.com/kyligence/xuanwu-log/pkg/schedule"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

var (
	testQuery      = "{job=\"fluent-bit\",app=\"yinglong\"}"
	testWorkingDir = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test"

	testBackup = &schedule.Backup{
		Data: &data.DataConf{
			Loki: &loki.LokiConf{
				Address: "http://aafdd592dddec49ed8bf3c35d9d538c9-577636166.us-west-2.elb.amazonaws.com:80",
			},
		},
		Queries: []*schedule.QueryConf{
			{
				Query: testQuery,
				Schedule: &schedule.Schedule{
					Interval: schedule.DefaultInterval,
					Max:      schedule.DefaultMax},
			},
		},
		Archive: &schedule.Archive{
			Type:        "zip",
			WorkingDir:  testWorkingDir,
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

func TestBackupReady(t *testing.T) {
	qc, ready := backupReady(testQuery, testBackup)
	if !ready {
		t.Fatalf("ready expected.")
	}

	t.Logf("qry conf: %+v", qc)
}

func TestGenerateRequests(t *testing.T) {
	startParsed, endParsed, _ := util.NormalizeTimes("now-6h", "now")

	queryConf, _ := backupReady(testQuery, testBackup)
	queryConf.Ensure(DOWNLOAD, testBackup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, nil, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	log.Printf("Requests total: %d", len(requests))
	for idx, request := range requests {
		log.Printf("Request #%d %s", idx+1, request.String())
	}
}

func TestGenerateRequestsHeadOnly(t *testing.T) {
	startParsed, endParsed, _ := util.NormalizeTimes("now-1h", "now")

	queryConf, _ := backupReady(testQuery, testBackup)
	queryConf.Ensure(DOWNLOAD, testBackup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, nil, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	log.Printf("Requests total: %d", len(requests))
	for idx, request := range requests {
		log.Printf("Request #%d %s", idx+1, request.String())
	}
}

func TestSubmit(t *testing.T) {
	startParsed, endParsed, _ := util.NormalizeTimes("now-6h", "now")

	queryConf, _ := backupReady(testQuery, testBackup)
	queryConf.Ensure(DOWNLOAD, testBackup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, testData, testStore)
	if err != nil {
		t.Fatalf("%s", err)
	}

	submit(requests)
}

func TestProceedNotExist(t *testing.T) {
	startParsed, endParsed, _ := util.NormalizeTimes("now-1h", "now")

	queryConf, _ := backupReady(testQuery, testBackup)
	queryConf.Ensure(DOWNLOAD, testBackup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, nil, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	dst := filepath.Join(testWorkingDir, "all.txt")
	err = proceed(dst, requests)
	if err != nil {
		t.Logf("expect error: %s", err)
	}
}

func TestCleanupNotExist(t *testing.T) {
	startParsed, endParsed, _ := util.NormalizeTimes("now-1h", "now")

	queryConf, _ := backupReady(testQuery, testBackup)
	queryConf.Ensure(DOWNLOAD, testBackup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, nil, nil)
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = cleanup(requests)
	if err != nil {
		t.Logf("expect error: %s", err)
	}
}
