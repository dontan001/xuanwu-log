package schedule

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/util"
)

func Run() {
	queryConf := &QueryConf{
		Query: "{job=\"fluent-bit\",app=\"yinglong\"}",
		Schedule: Schedule{
			Interval: DefaultInterval,
			Max:      DefaultMax},
		Prefix:      "",
		NamePattern: "log.%s.%s",
	}

	requests := generateRequests(queryConf)
	log.Printf("Requests total: %d", len(requests))

	submit(requests)
}

func generateRequests(conf *QueryConf) (requests []BackupRequest) {
	t := time.Now()
	log.Printf("Checked at: %s", t.Format(time.RFC3339Nano))

	lastBackup := util.CalcLastBackup(conf.Schedule.Interval, t)
	for idx := 1; idx <= conf.Schedule.Max; idx++ {
		start := lastBackup.Add(-time.Duration(conf.Schedule.Interval) * time.Hour)
		end := lastBackup.Add(-1 * time.Nanosecond)

		prefix := conf.Prefix
		if prefix == "" {
			prefix = fmt.Sprintf("%d", util.Hash(conf.Query))
		}
		archiveName := fmt.Sprintf(conf.NamePattern, prefix, lastBackup.Format(time.RFC3339))

		requests = append(requests, BackupRequest{
			Query: conf.Query,
			Start: start,
			End:   end,
			Archive: Archive{
				Name: archiveName,
			},
		})
		lastBackup = start
	}
	return
}

func submit(requests []BackupRequest) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, PARALLELIZE)
	realtime()

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
			log.Printf("proceed req: %s", req.String())
			req.Do()
			realtime()
		}(idx)
	}

	wg.Wait()
	realtime()
}

func realtime() {
	log.Printf("Routines total: %d", runtime.NumGoroutine())
}
