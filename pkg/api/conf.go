package api

import (
	"github.com/kyligence/xuanwu-log/pkg/data"
)

const (
	PARALLELIZE = 4
	trace       = false

	DOWNLOAD = "download"
)

type Server struct {
	Conf *ServerConf    `yaml:"server"`
	Data *data.DataConf `yaml:"data"`
}

type ServerConf struct {
	HttpPort   string `yaml:"port"`
	WorkingDir string `yaml:"workingDir"`
}
