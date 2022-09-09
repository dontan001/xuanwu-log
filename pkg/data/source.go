package data

import (
	"io"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data/loki/query"
)

func Extract(q string, start, end time.Time, result io.Writer) error {
	return query.QueryV2(q, start, end, result)
}
