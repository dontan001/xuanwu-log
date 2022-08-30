package schedule

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/log/loki/query"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

type BackupRequest struct {
	Query   string
	Start   time.Time
	End     time.Time
	Archive Archive
}

type Archive struct {
	// Prefix string
	Name string
}

func (req BackupRequest) Do() string {
	log.Printf("Proceed req: %s", req.String())

	fileName := req.Archive.Name
	fileNameZip := fmt.Sprintf("%s.zip", fileName)
	fileNameFull := fmt.Sprintf(BASE, fileName)
	fileNameZipFull := fmt.Sprintf(BASE, fileNameZip)
	result, e := os.Create(fileNameFull)
	if e != nil {
		log.Fatal(e)
	}

	defer func() {
		result.Close()
		if !trace {
			os.Remove(fileNameFull)
			os.Remove(fileNameZipFull)
		}
	}()
	query.QueryV2(req.Query, req.Start, req.End, result)

	if e := util.ZipSource(fileNameFull, fileNameZipFull); e != nil {
		log.Fatal(e)
	}

	return ""
}

func (req BackupRequest) String() string {
	var b strings.Builder

	b.WriteByte('{')
	b.WriteString("query=" + req.Query + ",")
	b.WriteString("start=" + req.Start.Format(time.RFC3339Nano) + ",")
	b.WriteString("startUnix=" + fmt.Sprintf("%d", req.Start.UnixNano()) + ",")
	b.WriteString("end=" + req.End.Format(time.RFC3339Nano) + ",")
	b.WriteString("endUnix=" + fmt.Sprintf("%d", req.End.UnixNano()))
	b.WriteByte('}')

	return b.String()
}
