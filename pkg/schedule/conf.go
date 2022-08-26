package schedule

import (
	"time"
)

var (
	intervalDefault = 6 * time.Hour
	maxDefault      = 3 * 4 * 6 * time.Hour
)

type Schedule struct {
	interval time.Duration
	max      time.Duration
}

type QueryConf struct {
	Query       string
	Schedule    Schedule
	Prefix      string
	NamePattern string
}
