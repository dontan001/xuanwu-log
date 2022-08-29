package schedule

import (
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
		Prefix:      "test",
		NamePattern: "log-%s",
	}

	requests := generateRequests(queryConf)
	log.Printf("Requests total: %d", len(requests))

	submit(requests)
}

func generateRequests(conf *QueryConf) (requests []QueryRequest) {
	t := time.Now()
	log.Printf("Checked at: %s", t.Format(time.RFC3339Nano))

	lastBackup := util.CalcLastBackup(conf.Schedule.Interval, t)
	for idx := 1; idx <= conf.Schedule.Max; idx++ {
		start := lastBackup.Add(-time.Duration(conf.Schedule.Interval) * time.Hour)
		end := lastBackup.Add(-1 * time.Nanosecond)

		requests = append(requests, QueryRequest{
			Query: conf.Query,
			Start: start,
			End:   end,
		})
		lastBackup = start
	}
	return
}

func submit(requests []QueryRequest) {
	var wg sync.WaitGroup
	ch := make(chan struct{}, PARALLELIZE)
	realtime()

	for idx := 0; idx < len(requests); idx++ {
		ch <- struct{}{}
		wg.Add(1)

		go func(i int) {
			defer func() {
				<-ch
				wg.Done()
			}()

			proceed(requests[i])
			realtime()
		}(idx)
	}

	wg.Wait()
	realtime()
}

func realtime() {
	log.Printf("Routines total: %d", runtime.NumGoroutine())
}