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
	http.HandleFunc("/log", func(w http.ResponseWriter, req *http.Request) {
		qry := req.URL.Query().Get("query")
		start := req.URL.Query().Get("start")
		end := req.URL.Query().Get("end")
		log.Printf("input: query=%s, start=%s, end=%s \n", qry, start, end)

		startParsed, endParsed, e := util.NormalizeTimes(start, end)
		if e != nil {
			log.Fatalf("time normalization err: %s", e)
		}
		log.Printf("parsed: query=%s, start=%s, end=%s", qry,
			startParsed.Format(time.RFC3339), endParsed.Format(time.RFC3339))

		buf := new(bytes.Buffer)
		writer := zip.NewWriter(buf)
		filename := "logs.all"
		f, err := writer.Create(filename)
		if err != nil {
			log.Fatal(err)
		}

		/*data := "1234567890"
		_, err = f.Write([]byte(data))*/
		result := &bytes.Buffer{}
		query.Query(qry, startParsed, endParsed, result)
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

	http.ListenAndServe(":8080", nil)
}
