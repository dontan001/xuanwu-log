package schedule

import (
	"fmt"
	"log"
	"os"

	"github.com/kyligence/xuanwu-log/pkg/log/loki/query"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

func proceed(req BackupRequest) string {
	log.Printf("proceed req: %s", req.String())

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
