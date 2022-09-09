package schedule

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

type BackupRequest struct {
	Query         string
	Start         time.Time
	End           time.Time
	ArchiveConfig ArchiveConfig
}

type ArchiveConfig struct {
	Name        string
	ArchiveName string
	WorkingDir  string

	ObjectPrefix string
}

func (req BackupRequest) Do() error {
	log.Printf("Proceed req: %s", req.String())

	exist, err := storage.Exist(req.ArchiveConfig.ObjectPrefix)
	if err != nil {
		return err
	}
	if exist {
		log.Printf("Object %q exists, skip", req.ArchiveConfig.ObjectPrefix)
	}

	fileName := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.Name)
	fileNameArchive := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.ArchiveName)

	result, err := os.Create(fileName)
	if err != nil {
		return err
	}

	defer func() {
		result.Close()
		if !trace {
			os.Remove(fileName)
			os.Remove(fileNameArchive)
		}
	}()

	err = data.Extract(req.Query, req.Start, req.End, result)
	if err != nil {
		return err
	}

	err = util.Compress(fileName, fileNameArchive)
	if err != nil {
		return err
	}

	if trace {
		return nil
	}
	err = storage.Upload(req.ArchiveConfig.ObjectPrefix, fileNameArchive)
	if err != nil {
		return err
	}

	return nil
}

func (req BackupRequest) String() string {
	var b strings.Builder

	b.WriteByte('{')
	b.WriteString("query=" + req.Query + ",")
	b.WriteString("start=" + req.Start.Format(time.RFC3339Nano) + ",")
	b.WriteString("startUnix=" + fmt.Sprintf("%d", req.Start.UnixNano()) + ",")
	b.WriteString("end=" + req.End.Format(time.RFC3339Nano) + ",")
	b.WriteString("endUnix=" + fmt.Sprintf("%d", req.End.UnixNano()) + ",")
	b.WriteString("prefix=" + req.ArchiveConfig.ObjectPrefix)
	b.WriteByte('}')

	return b.String()
}
