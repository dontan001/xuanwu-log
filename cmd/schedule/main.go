package main

import (
	"log"
	"os"

	"github.com/kyligence/xuanwu-log/pkg/schedule"
)

func main() {
	log.SetOutput(os.Stderr)

	schedule.Run()
}
