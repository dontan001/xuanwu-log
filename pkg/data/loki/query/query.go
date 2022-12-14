package query

import (
	"io"
	"log"
	"math"
	"time"

	cliClient "github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/logcli/output"
	"github.com/grafana/loki/pkg/logcli/query"

	"github.com/kyligence/xuanwu-log/pkg/data/loki/client"
	"github.com/kyligence/xuanwu-log/pkg/data/loki/query/v2"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

var (
	queryClient = client.NewQueryClient()

	mode      = "default"
	batchSize = 5000
	limit     = math.MaxInt
)

func Query(q string, start, end time.Time, result io.Writer) {
	defer util.TimeMeasure("query")()

	rangeQuery := newQuery(q, start, end)
	outputOptions := &output.LogOutputOptions{
		Timezone:      time.Local,
		NoLabels:      false,
		ColoredOutput: false,
	}

	out, err := output.NewLogOutput(result, mode, outputOptions)
	if err != nil {
		log.Fatalf("Unable to create log output: %s", err)
	}

	rangeQuery.DoQuery(queryClient, out, false)
}

func newQuery(q string, start, end time.Time) *query.Query {
	qry := &query.Query{}
	qry.QueryString = q

	qry.Start = start
	qry.End = end

	qry.BatchSize = batchSize
	qry.Limit = limit
	return qry
}

func QueryV2(client cliClient.Client, q string, start, end time.Time, result io.Writer) error {
	defer util.TimeMeasure("queryV2")()

	rangeQuery := newQueryV2(q, start, end)
	outputOptions := &output.LogOutputOptions{
		Timezone:      time.Local,
		NoLabels:      false,
		ColoredOutput: false,
	}

	/*out, err := output.NewLogOutput(result, mode, outputOptions)
	if err != nil {
		log.Printf("Unable to create log output: %s", err)
		return err
	}*/
	outV2 := &DefaultOutputV2{
		w:       result,
		options: outputOptions,
	}

	err := rangeQuery.DoQuery(client, outV2, false)
	if err != nil {
		log.Printf("DoQuery failed: %s", err)
		return err
	}

	return nil
}

func newQueryV2(q string, start, end time.Time) *v2.Query {
	qry := &v2.Query{}
	qry.QueryString = q

	qry.Start = start
	qry.End = end

	qry.BatchSize = batchSize
	qry.Limit = limit
	return qry
}
