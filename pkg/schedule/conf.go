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
	Data    *Data        `yaml:"data"`
	Queries []*QueryConf `yaml:"queries"`
	Archive *Archive     `yaml:"archive"`
}

type Data struct {
	Loki *Loki `yaml:"loki"`
}

type Loki struct {
	Address string `yaml:"address"`
}

func (c *Data) Validate() error {
	return nil
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
	Type        string `yaml:"type"`
	WorkingDir  string `yaml:"workingDir"`
	NamePattern string `yaml:"namePattern"`
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
