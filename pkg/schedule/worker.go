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

	Data  *data.Data
	Store *storage.Store
}

type ArchiveConfig struct {
	Name        string
	ArchiveName string
	WorkingDir  string

	ObjectPrefix string
}

func (req BackupRequest) Do() error {
	log.Printf("Proceed req: %s", req.String())

	exist, err := req.Store.Exist(req.ArchiveConfig.ObjectPrefix)
	if err != nil {
		return err
	}
	if exist {
		log.Printf("Backup %q exists, skip", req.ArchiveConfig.ObjectPrefix)
		return nil
	}

	fileName := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.Name)
	fileNameArchive := filepath.Join(req.ArchiveConfig.WorkingDir, req.ArchiveConfig.ArchiveName)

	defer func() {
		if !trace {
			err := os.Remove(fileName)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("Remove file %q with error: %s", fileName, err)
			}

			err = os.Remove(fileNameArchive)
			if err != nil && !os.IsNotExist(err) {
				log.Printf("Remove file %q with error: %s", fileNameArchive, err)
			}
		} else {
			log.Printf("trace mode, skips remove for %q", fileName)
		}
	}()

	err = req.Data.Extract(req.Query, req.Start, req.End, fileName)
	if err != nil {
		return err
	}

	err = util.Compress(fileName, fileNameArchive)
	if err != nil {
		return err
	}

	if trace {
		log.Printf("trace mode, skips upload for %q", req.ArchiveConfig.ObjectPrefix)
		return nil
	}
	err = req.Store.Upload(req.ArchiveConfig.ObjectPrefix, fileNameArchive)
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
