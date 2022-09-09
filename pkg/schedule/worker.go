package schedule

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
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

func (req BackupRequest) Do() error {
	log.Printf("Proceed req: %s", req.String())

	fileName := req.Archive.Name
	fileNameZip := fmt.Sprintf("%s.zip", fileName)
	fileNameFull := fmt.Sprintf(BASE, fileName)
	fileNameZipFull := fmt.Sprintf(BASE, fileNameZip)
	result, err := os.Create(fileNameFull)
	if err != nil {
		return err
	}

	defer func() {
		result.Close()
		if !trace {
			os.Remove(fileNameFull)
			os.Remove(fileNameZipFull)
		}
	}()
	err = data.Extract(req.Query, req.Start, req.End, result)
	if err != nil {
		return err
	}

	err = util.Compress(fileNameFull, fileNameZipFull)
	if err != nil {
		return err
	}

	/*err = s3.PutObject("", fileNameZipFull)
	if err != nil {
		return err
	}*/

	return nil
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
