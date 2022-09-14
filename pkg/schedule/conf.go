package schedule

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8

	DefaultType       = "zip"
	DefaultWorkingDir = "/var/log"
	// "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test"

	PARALLELIZE = 4
	trace       = true
)

type BackupConf struct {
	Queries []*QueryConf
	Archive *Archive
}

type QueryConf struct {
	Query    string
	Schedule *Schedule
	Archive  *ArchiveQuery
	Hash     string // sys
}

type Schedule struct {
	Interval int
	Max      int
}

type Archive struct {
	Type        string
	WorkingDir  string
	NamePattern string
}

type ArchiveQuery struct {
	Archive
	SubDir string // sys
}
