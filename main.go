package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/loki/query"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

func main() {
	log.SetOutput(os.Stderr)

	start()
}

func start() {
	http.HandleFunc("/log/big", func(w http.ResponseWriter, req *http.Request) {
		defer util.TimeMeasure("download")()

		fileName := req.URL.Query().Get("file")
		base := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/%s"
		fileNameFull := fmt.Sprintf(base, fileName)
		// fileNameFull := fmt.Sprintf("/app/test/%s", fileName)
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
		query.Query(qry, startParsed, endParsed, result)

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

	http.HandleFunc("/log/v2", func(w http.ResponseWriter, req *http.Request) {
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

		base := "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/%s"
		fileName := "tmp.txt"
		fileNameZip := fmt.Sprintf("%s.zip", fileName)
		fileNameFull := fmt.Sprintf(base, fileName)
		fileNameZipFull := fmt.Sprintf(base, fileNameZip)
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
		query.Query(qry, startParsed, endParsed, result)

		if e := util.ZipSource(fileNameFull, fileNameZipFull); e != nil {
			log.Fatal(e)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileNameZip))
		http.ServeFile(w, req, fileNameZipFull)
	})

	http.ListenAndServe(":8080", nil)
}
