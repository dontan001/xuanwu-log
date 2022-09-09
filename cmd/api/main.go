package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kyligence/xuanwu-log/pkg/data"
	"github.com/kyligence/xuanwu-log/pkg/util"
)

const (
	BASE = "/Users/dongge.tan/Dev/workspace/GOPATH/github.com/Kyligence/xuanwu-log/test/%s"
)

func main() {
	log.SetOutput(os.Stderr)

	start()
}

func start() {
	http.HandleFunc("/log/big", func(w http.ResponseWriter, req *http.Request) {
		defer util.TimeMeasure("download")()

		fileName := req.URL.Query().Get("file")
		fileNameFull := fmt.Sprintf(BASE, fileName)
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
		err := data.Extract(qry, startParsed, endParsed, result)
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

		fileName := "tmp.txt"
		fileNameZip := fmt.Sprintf("%s.zip", fileName)
		fileNameFull := fmt.Sprintf(BASE, fileName)
		fileNameZipFull := fmt.Sprintf(BASE, fileNameZip)
		result, err := os.Create(fileNameFull)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			result.Close()
			os.Remove(fileNameFull)
			os.Remove(fileNameZipFull)
		}()
		err = data.Extract(qry, startParsed, endParsed, result)
		if err != nil {
			log.Fatal(err)
		}

		if err = util.ZipSource(fileNameFull, fileNameZipFull); err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/zip")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileNameZip))
		http.ServeFile(w, req, fileNameZipFull)
	})

	http.ListenAndServe(":8080", nil)
}
