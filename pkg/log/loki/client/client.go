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
	client.Address = "http://a82482b9a9c354066bebaae7008def97-1902638330.us-west-2.elb.amazonaws.com:80"

	return client
}
