package schedule

import (
	"log"
)

func proceed(req QueryRequest) string {

	log.Printf("proceed req: %s", req.String())
	return ""
}
