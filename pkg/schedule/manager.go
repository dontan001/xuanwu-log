package schedule

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

func Run(backup *Backup) {
	/*backup := &Backup{
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
	}*/

	data := func(c *Backup) *data.Data {
		d := &data.Data{Conf: c.Data}
		d.Setup()
		return d
	}(backup)

	store := func(c *Backup) *storage.Store {
		s := &storage.Store{Config: c.Archive.S3}
		s.Setup()
		return s
	}(backup)

	for _, queryConf := range backup.Queries {
		log.Printf("Proceed qry: %q", queryConf.Query)
		queryConf.Ensure(BACKUP, backup)
		requests := queryConf.generateRequests(data, store)

		log.Printf("Requests total: %d", len(requests))
		submit(requests)
	}
}

func (conf *QueryConf) Ensure(sub string, backup *Backup) {
	if conf.Archive == nil {
		conf.Archive = &ArchiveQuery{}
	}
	if conf.Archive.WorkingDir == "" {
		conf.Archive.WorkingDir = backup.Archive.WorkingDir
	}
	if conf.Archive.Type == "" {
		conf.Archive.Type = backup.Archive.Type
	}
	if conf.Archive.NamePattern == "" {
		conf.Archive.NamePattern = backup.Archive.NamePattern
	}

	conf.Hash = fmt.Sprintf("%d", util.Hash(conf.Query))
	conf.Archive.SubDir = filepath.Join(sub, conf.Hash)

	fDir := filepath.Join(conf.Archive.WorkingDir, conf.Archive.SubDir)
	err := os.MkdirAll(fDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (conf *QueryConf) generateRequests(data *data.Data, store *storage.Store) (requests []BackupRequest) {
	t := time.Now()
	log.Printf("Checked at: %s", t.Format(time.RFC3339Nano))

	lastBackup := util.CalcLastBackup(conf.Schedule.Interval, t)
	for idx := 1; idx <= conf.Schedule.Max; idx++ {
		start := lastBackup.Add(-time.Duration(conf.Schedule.Interval) * time.Hour)
		end := lastBackup.Add(-1 * time.Nanosecond)

		name := fmt.Sprintf(conf.Archive.NamePattern, lastBackup.Format(time.RFC3339))
		archiveName := fmt.Sprintf("%s.%s", name, conf.Archive.Type)

		requests = append(requests, BackupRequest{
			Query: conf.Query,
			Start: start,
			End:   end,
			ArchiveConfig: ArchiveConfig{
				Name:         name,
				ArchiveName:  archiveName,
				WorkingDir:   filepath.Join(conf.Archive.WorkingDir, conf.Archive.SubDir),
				ObjectPrefix: filepath.Join(conf.Archive.SubDir, archiveName),
			},

			Data:  data,
			Store: store,
		})
		lastBackup = start
	}
	return
}

func submit(requests []BackupRequest) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, PARALLELIZE)
	realtime()

	log.Printf("Request total: %d", len(requests))
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
				log.Printf("Retry next time...")
			}
			realtime()
		}(idx)
	}

	wg.Wait()
	realtime()
}

func realtime() {
	if !trace {
		return
	}

	log.Printf("Routines total: %d", runtime.NumGoroutine())
}
