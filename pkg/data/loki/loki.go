package loki

import (
	"io"
	"time"

	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/kyligence/xuanwu-log/pkg/data/loki/query"
	"github.com/prometheus/common/config"
)

type LokiConf struct {
	Address string `yaml:"address"`
}

type Loki struct {
	Conf   *LokiConf
	client client.Client
}

func (loki *Loki) Setup() {
	if loki.Conf != nil {
		client := &client.DefaultClient{
			TLSConfig: config.TLSConfig{},
		}

		// client.Address = "http://loki-loki-distributed-gateway.loki:80"
		// client.Address = "http://aafdd592dddec49ed8bf3c35d9d538c9-577636166.us-west-2.elb.amazonaws.com:80"
		client.Address = loki.Conf.Address

		loki.client = client
	}
}

func (loki *Loki) Query(q string, start, end time.Time, result io.Writer) error {
	return query.QueryV2(loki.client, q, start, end, result)
}
