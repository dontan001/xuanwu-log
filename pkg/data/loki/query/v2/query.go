package v2

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/fatih/color"
	"github.com/grafana/loki/pkg/logcli/client"
	"github.com/grafana/loki/pkg/logcli/output"
	"github.com/grafana/loki/pkg/logcli/query"
	"github.com/grafana/loki/pkg/loghttp"
	"github.com/grafana/loki/pkg/logproto"
	"github.com/grafana/loki/pkg/logql"
	"github.com/grafana/loki/pkg/logql/stats"
)

type Query struct {
	query.Query
}

type streamEntryPair struct {
	entry  loghttp.Entry
	labels loghttp.LabelSet
}

func (q Query) DoQuery(c client.Client, out output.LogOutput, statistics bool) error {
	d := q.resultsDirection()

	if q.Limit < q.BatchSize {
		q.BatchSize = q.Limit
	}
	resultLength := 0
	total := 0
	start := q.Start
	end := q.End
	var lastEntry []*loghttp.Entry
	for total < q.Limit {
		bs := q.BatchSize
		// We want to truncate the batch size if the remaining number
		// of items needed to reach the limit is less than the batch size
		if q.Limit-total < q.BatchSize {
			// Truncated batchsize is q.Limit - total, however we add to this
			// the length of the overlap from the last query to make sure we get the
			// correct amount of new logs knowing there will be some overlapping logs returned.
			bs = q.Limit - total + len(lastEntry)
		}
		resp, err := c.QueryRange(q.QueryString, bs, start, end, d, q.Step, q.Interval, q.Quiet)
		if err != nil {
			// log.Fatalf("Query failed: %+v", err)
			err := fmt.Errorf("QueryRange failed: %+v", err)
			return err
		}

		if statistics {
			q.printStats(resp.Data.Statistics)
		}

		resultLength, lastEntry = q.printResult(resp.Data.Result, out, lastEntry)
		// Was not a log stream query, or no results, no more batching
		if resultLength <= 0 {
			break
		}
		// Also no result, wouldn't expect to hit this.
		if len(lastEntry) == 0 {
			break
		}
		// Can only happen if all the results return in one request
		if resultLength == q.Limit {
			break
		}
		if len(lastEntry) >= q.BatchSize {
			err := fmt.Errorf("Invalid batch size %v, the next query will have %v overlapping entries "+
				"(there will always be 1 overlapping entry but Loki allows multiple entries to have "+
				"the same timestamp, so when a batch ends in this scenario the next query will include "+
				"all the overlapping entries again).  Please increase your batch size to at least %v to account "+
				"for overlapping entryes\n", q.BatchSize, len(lastEntry), len(lastEntry)+1)
			return err
		}

		// Batching works by taking the timestamp of the last query and using it in the next query,
		// because Loki supports multiple entries with the same timestamp it's possible for a batch to have
		// fallen in the middle of a list of entries for the same time, so to make sure we get all entries
		// we start the query on the same time as the last entry from the last batch, and then we keep this last
		// entry and remove the duplicate when printing the results.
		// Because of this duplicate entry, we have to subtract it here from the total for each batch
		// to get the desired limit.
		total += resultLength
		// Based on the query direction we either set the start or end for the next query.
		// If there are multiple entries in `lastEntry` they have to have the same timestamp so we can pick just the first
		if q.Forward {
			start = lastEntry[0].Timestamp
		} else {
			// The end timestamp is exclusive on a backward query, so to make sure we get back an overlapping result
			// fudge the timestamp forward in time to make sure to get the last entry from this batch in the next query
			end = lastEntry[0].Timestamp.Add(1 * time.Nanosecond)
		}

	}

	return nil
}

func (q *Query) printResult(value loghttp.ResultValue, out output.LogOutput, lastEntry []*loghttp.Entry) (int, []*loghttp.Entry) {
	length := -1
	var entry []*loghttp.Entry
	switch value.Type() {
	case logql.ValueTypeStreams:
		length, entry = q.printStream(value.(loghttp.Streams), out, lastEntry)
	default:
		log.Fatalf("Unable to print unsupported type: %v", value.Type())
	}
	return length, entry
}

func (q *Query) printStream(streams loghttp.Streams, out output.LogOutput, lastEntry []*loghttp.Entry) (int, []*loghttp.Entry) {
	common := commonLabels(streams)

	// Remove the labels we want to show from common
	if len(q.ShowLabelsKey) > 0 {
		common = matchLabels(false, common, q.ShowLabelsKey)
	}

	if len(common) > 0 && !q.Quiet {
		log.Println("Common labels:", color.RedString(common.String()))
	}

	if len(q.IgnoreLabelsKey) > 0 && !q.Quiet {
		log.Println("Ignoring labels key:", color.RedString(strings.Join(q.IgnoreLabelsKey, ",")))
	}

	if len(q.ShowLabelsKey) > 0 && !q.Quiet {
		log.Println("Print only labels key:", color.RedString(strings.Join(q.ShowLabelsKey, ",")))
	}

	// Remove ignored and common labels from the cached labels and
	// calculate the max labels length
	maxLabelsLen := q.FixedLabelsLen
	for i, s := range streams {
		// Remove common labels
		ls := subtract(s.Labels, common)

		if len(q.ShowLabelsKey) > 0 {
			ls = matchLabels(true, ls, q.ShowLabelsKey)
		}

		// Remove ignored labels
		if len(q.IgnoreLabelsKey) > 0 {
			ls = matchLabels(false, ls, q.IgnoreLabelsKey)
		}

		// Overwrite existing Labels
		streams[i].Labels = ls

		// Update max labels length
		len := len(ls.String())
		if maxLabelsLen < len {
			maxLabelsLen = len
		}
	}

	// sort and display entries
	allEntries := make([]streamEntryPair, 0)

	for _, s := range streams {
		for _, e := range s.Entries {
			allEntries = append(allEntries, streamEntryPair{
				entry:  e,
				labels: s.Labels,
			})
		}
	}

	if len(allEntries) == 0 {
		return 0, nil
	}

	if q.Forward {
		sort.Slice(allEntries, func(i, j int) bool { return allEntries[i].entry.Timestamp.Before(allEntries[j].entry.Timestamp) })
	} else {
		sort.Slice(allEntries, func(i, j int) bool { return allEntries[i].entry.Timestamp.After(allEntries[j].entry.Timestamp) })
	}

	printed := 0
	for _, e := range allEntries {
		// Skip the last entry if it overlaps, this happens because batching includes the last entry from the last batch
		if len(lastEntry) > 0 && e.entry.Timestamp == lastEntry[0].Timestamp {
			skip := false
			// Because many logs can share a timestamp in the unlucky event a batch ends with a timestamp
			// shared by multiple entries we have to check all that were stored to see if we've already
			// printed them.
			for _, le := range lastEntry {
				if e.entry.Line == le.Line {
					skip = true
				}
			}
			if skip {
				continue
			}
		}
		out.FormatAndPrintln(e.entry.Timestamp, e.labels, maxLabelsLen, e.entry.Line)
		printed++
	}

	// Loki allows multiple entries at the same timestamp, this is a bit of a mess if a batch ends
	// with an entry that shared multiple timestamps, so we need to keep a list of all these entries
	// because the next query is going to contain them too and we want to not duplicate anything already
	// printed.
	lel := []*loghttp.Entry{}
	// Start with the timestamp of the last entry
	le := allEntries[len(allEntries)-1].entry
	for i, e := range allEntries {
		// Save any entry which has this timestamp (most of the time this will only be the single last entry)
		if e.entry.Timestamp.Equal(le.Timestamp) {
			lel = append(lel, &allEntries[i].entry)
		}
	}

	return printed, lel
}

func (q *Query) printStats(stats stats.Result) {
	writer := tabwriter.NewWriter(os.Stderr, 0, 8, 0, '\t', 0)
	stats.Log(kvLogger{Writer: writer})
}

func (q *Query) resultsDirection() logproto.Direction {
	if q.Forward {
		return logproto.FORWARD
	}
	return logproto.BACKWARD
}

type kvLogger struct {
	*tabwriter.Writer
}

func (k kvLogger) Log(keyvals ...interface{}) error {
	for i := 0; i < len(keyvals); i += 2 {
		fmt.Fprintln(k.Writer, color.BlueString("%s", keyvals[i]), "\t", fmt.Sprintf("%v", keyvals[i+1]))
	}
	k.Flush()
	return nil
}
