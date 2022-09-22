package query

import (
	"fmt"
	"hash/fnv"
	"io"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/grafana/loki/pkg/logcli/output"
	"github.com/grafana/loki/pkg/loghttp"
)

type DefaultOutputV2 struct {
	w       io.Writer
	options *output.LogOutputOptions
}

func (o *DefaultOutputV2) FormatAndPrintln(ts time.Time, lbls loghttp.LabelSet, maxLabelsLen int, line string) {
	timestamp := ts.In(o.options.Timezone).Format(time.RFC3339)
	line = strings.TrimSpace(line)

	if o.options.NoLabels {
		fmt.Fprintf(o.w, "%s %s\n", color.BlueString(timestamp), line)
		return
	}
	if o.options.ColoredOutput {
		labelsColor := getColor(lbls.String()).SprintFunc()
		fmt.Fprintf(o.w, "%s %s %s\n", color.BlueString(timestamp), labelsColor(padLabel(lbls, maxLabelsLen)), line)
	} else {
		fmt.Fprintf(o.w, "%s %s %s\n", timestamp, padLabel(lbls, maxLabelsLen), line)
	}

}

func padLabel(ls loghttp.LabelSet, maxLabelsLen int) string {
	labels := ls.String()
	if len(labels) < maxLabelsLen {
		labels += strings.Repeat(" ", maxLabelsLen-len(labels))
	}
	return labels
}

func getColor(labels string) *color.Color {
	hash := fnv.New32()
	_, _ = hash.Write([]byte(labels))
	id := hash.Sum32() % uint32(len(colorList))
	color := colorList[id]
	return color
}

var colorList = []*color.Color{
	color.New(color.FgHiCyan),
	color.New(color.FgCyan),
	color.New(color.FgHiGreen),
	color.New(color.FgGreen),
	color.New(color.FgHiMagenta),
	color.New(color.FgMagenta),
	color.New(color.FgHiYellow),
	color.New(color.FgYellow),
	color.New(color.FgHiRed),
	color.New(color.FgRed),
}
