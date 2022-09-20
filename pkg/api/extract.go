package api

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/schedule"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

type ExtractRequest struct {
	Query string
	Start time.Time
	End   time.Time

	ArchiveConfig schedule.ArchiveConfig

	FromData bool
	Data     *data.Data
	Store    *storage.Store
}

func backupReady(q string, backup *schedule.Backup) (*schedule.QueryConf, bool) {
	if backup == nil {
		return nil, false
	}

	for _, qry := range backup.Queries {
		if qry.Query == q {
			return qry, true
		}
	}

	return nil, false
}

func extractWithBackup(startParsed time.Time, endParsed time.Time,
	queryConf *schedule.QueryConf,
	backup *schedule.Backup,
	data *data.Data,
	store *storage.Store) error {

	queryConf.Ensure(DOWNLOAD, backup)
	requests, err := generateRequests(startParsed, endParsed, queryConf, data, store)
	if err != nil {
		return err
	}

	submit(requests)
	proceed(requests)
	cleanup(requests)

	return nil
}

func generateRequests(start, end time.Time,
	conf *schedule.QueryConf,
	data *data.Data,
	store *storage.Store) ([]ExtractRequest, error) {

	t := time.Now()
	log.Printf("Checked at: %s", t.Format(time.RFC3339Nano))

	lastBackup := util.CalcLastBackup(conf.Schedule.Interval, t)
	if end.Before(lastBackup) {
		return nil, fmt.Errorf("end time is before backup last time")
	}

	var requests []ExtractRequest
	for idx := 1; idx <= conf.Schedule.Max; idx++ {
		startTime := lastBackup.Add(-time.Duration(conf.Schedule.Interval) * time.Hour)
		endTime := lastBackup.Add(-1 * time.Nanosecond)

		if start.After(endTime) {
			break
		}

		headRequest := false
		if start.After(startTime) {
			headRequest = true
		}

		startTimeRequest := startTime
		if headRequest {
			startTimeRequest = start
		}

		name := fmt.Sprintf(conf.Archive.NamePattern, lastBackup.Format(time.RFC3339))
		archiveName := fmt.Sprintf("%s.%s", name, conf.Archive.Type)

		requests = append(requests, ExtractRequest{
			Query: conf.Query,
			Start: startTimeRequest,
			End:   endTime,
			ArchiveConfig: schedule.ArchiveConfig{
				Name:         name,
				ArchiveName:  archiveName,
				WorkingDir:   filepath.Join(conf.Archive.WorkingDir, conf.Archive.SubDir),
				ObjectPrefix: filepath.Join(conf.Archive.SubDir, archiveName),
			},

			FromData: headRequest,
			Data:     data,
			Store:    store,
		})
		lastBackup = startTime
	}

	startTimeRequest := util.CalcLastBackup(conf.Schedule.Interval, t)
	tailRequest := ExtractRequest{
		Query: conf.Query,
		Start: startTimeRequest,
		End:   end,
		ArchiveConfig: schedule.ArchiveConfig{
			Name:        "tail.log",
			ArchiveName: "tail.log.zip",
			WorkingDir:  filepath.Join(conf.Archive.WorkingDir, conf.Archive.SubDir),
		},

		FromData: true,
		Data:     data,
		Store:    store,
	}

	requests = append(requests, tailRequest)
	return requests, nil
}

func submit(requests []ExtractRequest) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, PARALLELIZE)

	for idx := 0; idx < len(requests); idx++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(i int) {
			defer func() {
				if r := recover(); r != nil {
					log.Fatalf("Submit with err: %s", r)
				}
				<-ch
				wg.Done()
			}()

			req := requests[i]
			err := req.Do()
			if err != nil {
				log.Printf("Proceed req err: %s", err)
			}
		}(idx)
	}

	wg.Wait()
}

func proceed(requests []ExtractRequest) {

}

func cleanup(requests []ExtractRequest) {

}

func (req ExtractRequest) Do() error {
	log.Printf("Proceed extract req: %s", req.String())

	fileName := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.Name)
	fileNameArchive := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.ArchiveName)

	if req.FromData {
		log.Printf("Extract from Loki for req: %s", req.String())
	} else {
		log.Printf("Extract from remote storage for req: %s", req.String())
		return req.Store.Download(req.ArchiveConfig.ObjectPrefix, fileNameArchive)
	}

	result, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		result.Close()
	}()

	err = req.Data.Extract(req.Query, req.Start, req.End, result)
	if err != nil {
		return err
	}

	return nil
}

func (req ExtractRequest) String() string {
	var b strings.Builder

	b.WriteByte('{')
	b.WriteString("query=" + req.Query + ",")
	b.WriteString("start=" + req.Start.Format(time.RFC3339Nano) + ",")
	b.WriteString("startUnix=" + fmt.Sprintf("%d", req.Start.UnixNano()) + ",")
	b.WriteString("end=" + req.End.Format(time.RFC3339Nano) + ",")
	b.WriteString("endUnix=" + fmt.Sprintf("%d", req.End.UnixNano()) + ",")
	b.WriteString("data=" + fmt.Sprintf("%t", req.FromData) + ",")
	b.WriteString("prefix=" + req.ArchiveConfig.ObjectPrefix)
	b.WriteByte('}')

	return b.String()
}
