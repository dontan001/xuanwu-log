package api

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/schedule"
	"github.com/kyligence/xuanwu-log/pkg/storage"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

func Start(server *Server, backup *schedule.Backup) {
	data := func(s *Server) *data.Data {
		d := &data.Data{Conf: s.Data}
		d.Setup()
		return d
	}(server)

	store := func(c *schedule.Backup) *storage.Store {
		s := &storage.Store{Config: c.Archive.S3}
		s.Setup()
		return s
	}(backup)

	http.HandleFunc("/log/big", func(w http.ResponseWriter, req *http.Request) {
		defer util.TimeMeasure("download")()

		fileName := req.URL.Query().Get("file")
		fileNameFull := filepath.Join(server.Conf.WorkingDir, fileName)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		http.ServeFile(w, req, fileNameFull)
	})

	// io.Pipe
	http.HandleFunc("/log/pipe", func(w http.ResponseWriter, req *http.Request) {
	})

	// transfer chunk
	http.HandleFunc("/log/chunk", func(w http.ResponseWriter, req *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("expected http.ResponseWriter to be an http.Flusher")
		}
		w.Header().Set("X-Content-Type-Options", "nosniff")
		for i := 1; i <= 10; i++ {
			fmt.Fprintf(w, "Chunk #%d\n", i)
			flusher.Flush() // Trigger "chunked" encoding and send a chunk...
			time.Sleep(500 * time.Millisecond)
		}
	})

	http.HandleFunc("/log", func(w http.ResponseWriter, req *http.Request) {
		defer util.TimeMeasure("download")()

		qry := req.URL.Query().Get("query")
		start := req.URL.Query().Get("start")
		end := req.URL.Query().Get("end")
		log.Printf("input: query=%s, start=%s, end=%s \n", qry, start, end)

		startParsed, endParsed, e := util.NormalizeTimes(start, end)
		if e != nil {
			log.Fatalf("time normalization err: %s", e)
		}
		log.Printf("parsed: query=%s, start=%s [ %d ], end=%s [ %d ]", qry,
			startParsed.Format(time.RFC3339Nano), startParsed.UnixNano(), endParsed.Format(time.RFC3339Nano), endParsed.UnixNano())

		result := &bytes.Buffer{}
		err := data.ExtractWithWriter(qry, startParsed, endParsed, result)
		if err != nil {
			log.Fatal(err)
		}

		buf := new(bytes.Buffer)
		writer := zip.NewWriter(buf)
		filename := "logs.all"
		f, err := writer.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write(result.Bytes())
		if err != nil {
			log.Fatal(err)
		}
		err = writer.Close()
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.zip\"", filename))
		w.Write(buf.Bytes())
	})

	http.HandleFunc("/log/download", func(w http.ResponseWriter, req *http.Request) {
		defer util.TimeMeasure("download")()

		qry := req.URL.Query().Get("query")
		start := req.URL.Query().Get("start")
		end := req.URL.Query().Get("end")
		log.Printf("input: query=%s, start=%s, end=%s \n", qry, start, end)

		startParsed, endParsed, err := util.NormalizeTimes(start, end)
		if err != nil {
			log.Fatalf("time normalization err: %s", err)
		}
		log.Printf("parsed: query=%s, start=%s [ %d ], end=%s [ %d ]", qry,
			startParsed.Format(time.RFC3339Nano), startParsed.UnixNano(), endParsed.Format(time.RFC3339Nano), endParsed.UnixNano())

		fileName := util.RandomName()
		fileNameFull := filepath.Join(server.Conf.WorkingDir, fileName)
		fileNameArchive := fmt.Sprintf("%s.zip", fileName)
		fileNameArchiveFull := filepath.Join(server.Conf.WorkingDir, fileNameArchive)

		defer func() {
			if !trace {
				err := os.Remove(fileNameFull)
				if err != nil && !os.IsNotExist(err) {
					log.Printf("Remove file %q with error: %s", fileNameFull, err)
				}

				err = os.Remove(fileNameArchiveFull)
				if err != nil && !os.IsNotExist(err) {
					log.Printf("Remove file %q with error: %s", fileNameArchiveFull, err)
				}
			}
		}()

		queryConf, ready := backupReady(qry, backup)
		if !ready {
			err = data.Extract(qry, startParsed, endParsed, fileNameFull)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			err = extractWithBackup(startParsed, endParsed, queryConf, backup, fileNameFull, data, store)
			if err != nil {
				log.Fatal(err)
			}
		}

		if err = util.Compress(fileNameFull, fileNameArchiveFull); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileNameArchive))
		http.ServeFile(w, req, fileNameArchiveFull)
	})

	http.ListenAndServe(server.Conf.HttpPort, nil)
}
