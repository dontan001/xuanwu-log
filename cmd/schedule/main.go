package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kyligence/xuanwu-log/pkg/schedule"
	yaml "gopkg.in/yaml.v2"
)

var configFile = flag.String("config.file", "/etc/config/config.yaml", "backup config options")

func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()

	var config *schedule.Backup
	config, err := loadConf(*configFile)
	if err != nil {
		panic(err)
	}

	err = config.Validate()
	if err != nil {
		panic(fmt.Errorf("config validation err: %s", err))
	}

	err = os.MkdirAll(config.Archive.WorkingDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	schedule.Run(config)
}

func loadConf(fileName string) (*schedule.Backup, error) {
	log.Printf("Load config from %q", fileName)
	bytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var conf schedule.Backup
	err = yaml.Unmarshal(bytes, &conf)
	if err != nil {
		return nil, fmt.Errorf("parse config file %q error: %v", fileName, err)
	}

	return &conf, nil
}
