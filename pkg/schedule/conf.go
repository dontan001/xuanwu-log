package schedule

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8

	PARALLELIZE = 4
	BASE        = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/backup/%s"

	trace = true
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
