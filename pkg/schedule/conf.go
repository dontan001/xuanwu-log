package schedule

import (
	"strings"
	"time"
)

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8

	PARALLELIZE = 4
	BASE        = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/backup/%s"
)

type QueryConf struct {
	Query       string
	Schedule    Schedule
	Prefix      string
	NamePattern string
}

type Schedule struct {
	Interval int
	Max      int
}

type BackupRequest struct {
	Query   string
	Start   time.Time
	End     time.Time
	Archive Archive
}

type Archive struct {
	// Prefix string
	Name string
}

func (r *BackupRequest) String() string {
	var b strings.Builder

	b.WriteByte('{')
	b.WriteString("query=" + r.Query + ",")
	b.WriteString("start=" + r.Start.Format(time.RFC3339Nano) + ",")
	b.WriteString("end=" + r.End.Format(time.RFC3339Nano))
	b.WriteByte('}')

	return b.String()
}
