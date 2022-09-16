package schedule

import (
	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/storage/s3"
)

const (
	DefaultInterval = 3
	DefaultMax      = 1 * 8

	DefaultType       = "zip"
	DefaultWorkingDir = "/var/log"

	PARALLELIZE = 4
	trace       = true
)

type BackupConf struct {
	Data    *data.DataConf `yaml:"data"`
	Queries []*QueryConf   `yaml:"queries"`
	Archive *Archive       `yaml:"archive"`
}

type QueryConf struct {
	Query    string        `yaml:"query"`
	Schedule *Schedule     `yaml:"schedule"`
	Archive  *ArchiveQuery `yaml:"archive"`
	Hash     string        // sys
}

type Schedule struct {
	Interval int `yaml:"interval"`
	Max      int `yaml:"max"`
}

type Archive struct {
	Type        string       `yaml:"type"`
	WorkingDir  string       `yaml:"workingDir"`
	NamePattern string       `yaml:"namePattern"`
	S3          *s3.S3Config `yaml:"s3"`
}

type ArchiveQuery struct {
	Archive
	SubDir string // sys
}

func (c *QueryConf) Validate() error {
	return nil
}

func (c *BackupConf) Validate() error {
	return nil
}
