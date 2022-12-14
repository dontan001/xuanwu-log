package data

import (
	"io"
	"log"
	"os"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data/loki"
)

type DataConf struct {
	Loki *loki.LokiConf `yaml:"loki"`
}

func (dataConf *DataConf) Validate() error {
	return nil
}

type Data struct {
	Conf *DataConf
	Loki *loki.Loki
}

func (data *Data) Setup() {
	if data.Conf != nil {
		data.Loki = &loki.Loki{
			Conf: data.Conf.Loki,
		}

		data.Loki.Setup()
		log.Printf("Data setup finish.")
	}
}

func (data *Data) Extract(q string, start, end time.Time, fileName string) error {
	result, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		result.Close()
	}()

	return data.Loki.Query(q, start, end, result)
}

func (data *Data) ExtractWithWriter(q string, start, end time.Time, result io.Writer) error {
	return data.Loki.Query(q, start, end, result)
}
