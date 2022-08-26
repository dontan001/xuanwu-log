package schedule

import (
	"strings"
	"time"
)

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8
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

type QueryRequest struct {
	Query string
	Start time.Time
	End   time.Time
}

func (r *QueryRequest) String() string {
	var b strings.Builder

	b.WriteByte('{')
	b.WriteString("query=" + r.Query + ",")
	b.WriteString("start=" + r.Start.Format(time.RFC3339Nano) + ",")
	b.WriteString("end=" + r.End.Format(time.RFC3339Nano))
	b.WriteByte('}')

	return b.String()
}
