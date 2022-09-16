package api

import (
	"github.com/kyligence/xuanwu-log/pkg/data"
)

type Server struct {
	Conf *ServerConf    `yaml:"server"`
	Data *data.DataConf `yaml:"data"`
}

type ServerConf struct {
	HttpPort   string `yaml:"port"`
	WorkingDir string `yaml:"workingDir"`
}
