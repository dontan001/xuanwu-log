package client

import (
	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/prometheus/common/config"
)

func NewQueryClient() client.Client {
	client := &client.DefaultClient{
		TLSConfig: config.TLSConfig{},
	}

	// client.Address = "http://loki-loki-distributed-gateway.loki:80"
	client.Address = "http://aafdd592dddec49ed8bf3c35d9d538c9-577636166.us-west-2.elb.amazonaws.com:80"

	return client
}
