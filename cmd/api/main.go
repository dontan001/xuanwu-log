package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/kyligence/xuanwu-log/pkg/api"
	"github.com/kyligence/xuanwu-log/pkg/data"
	yaml "gopkg.in/yaml.v2"
)

var configFile = flag.String("config.file", "/etc/config/config.yaml", "api config options")

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

	data := func(s *api.Server) *data.Data {
		d := &data.Data{Conf: s.Data}
		d.Setup()
		return d
	}(server)

	api.Start(server, data)
	log.Printf("Log API server started.")
}
