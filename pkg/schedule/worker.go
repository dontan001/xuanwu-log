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
	fileName := req.Archive.Name
	fileNameZip := fmt.Sprintf("%s.zip", fileName)
	fileNameFull := fmt.Sprintf(BASE, fileName)
	fileNameZipFull := fmt.Sprintf(BASE, fileNameZip)
	result, e := os.Create(fileNameFull)
	if e != nil {
		log.Fatal(e)
	}

	defer func() {
		log.Printf("Clean up\n")
		result.Close()
		os.Remove(fileNameFull)
		os.Remove(fileNameZipFull)
	}()
	query.Query(req.Query, req.Start, req.End, result)

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
	b.WriteString("end=" + req.End.Format(time.RFC3339Nano))
	b.WriteByte('}')

	return b.String()
}
