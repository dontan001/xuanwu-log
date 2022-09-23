package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kyligence/xuanwu-log/pkg/api"
	"github.com/kyligence/xuanwu-log/pkg/schedule"
	yaml "gopkg.in/yaml.v2"
)

var configFile = flag.String("config.file", "/etc/config/config.yaml", "api config options")
var backupFile = flag.String("config.backup", "", "backup config options")

func main() {
	log.SetOutput(os.Stderr)
	flag.Parse()

	var server *api.Server
	server, err := func(fileName string) (*api.Server, error) {
		log.Printf("Load config from %q", fileName)
		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		var conf api.Server
		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			return nil, fmt.Errorf("parse config file %q error: %v", fileName, err)
		}

		return &conf, nil
	}(*configFile)
	if err != nil {
		panic(err)
	}

	backup, err := func(fileName string) (*schedule.Backup, error) {
		if fileName == "" {
			return nil, nil
		}

		log.Printf("Load backup from %q", fileName)
		bytes, err := ioutil.ReadFile(fileName)
		if err != nil {
			return nil, err
		}

		var conf schedule.Backup
		err = yaml.Unmarshal(bytes, &conf)
		if err != nil {
			return nil, fmt.Errorf("parse backup file %q error: %v", fileName, err)
		}

		return &conf, nil
	}(*backupFile)
	if err != nil {
		panic(err)
	}

	err = os.MkdirAll(server.Conf.WorkingDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	log.Printf("Log API server starting...")
	api.Start(server, backup)
}
