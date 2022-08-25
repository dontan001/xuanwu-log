package schedule

import "time"

var (
	interval = 6 * time.Hour
	max      = 3 * 4 * 6 * time.Hour
)

type schedule struct {
	interval time.Duration
	max      time.Duration
}

type queryConf struct {
	Query       string
	Schedule    schedule
	Prefix      string
	NamePattern string
}
