package schedule

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8

	DefaultType       = "zip"
	DefaultWorkingDir = "/var/log"
	// "/Users/dongge.tan/Dev/wozrkspace/GOPATH/github.com/Kyligence/xuanwu-log/test"

	PARALLELIZE = 4
	trace       = true
)

type QueryConf struct {
	Query    string
	Schedule *Schedule
	Archive  *Archive
	Hash     string // sys
}

type Schedule struct {
	Interval int
	Max      int
}

type Archive struct {
	Type        string
	WorkingDir  string
	SubDir      string // sys
	NamePattern string
}
