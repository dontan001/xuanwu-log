package util

import (
	"fmt"
	"log"
	"math"
	"regexp"
	"strconv"
	"time"
)

func NormalizeTimes(start, end string) (time.Time, time.Time, error) {
	startParsed, e := ParseTime(start)
	if e != nil {
		return time.Time{}, time.Time{}, e
	}

	endParsed, e := ParseTime(end)
	if e != nil {
		return startParsed, time.Time{}, e
	}

	return startParsed, endParsed, nil
}

func ParseTime(timeStr string) (time.Time, error) {
	v, err := strconv.ParseInt(timeStr, 10, 64)
	if err == nil {
		return time.Unix(0, v), nil
	}

	re, _ := regexp.Compile("^(now)([+-]\\d.*)?$")
	match := re.MatchString(timeStr)
	if !match {
		return time.Time{}, fmt.Errorf("time format not support")
	}

	subs := re.FindStringSubmatch(timeStr)
	if subs[2] == "" {
		return time.Now(), nil
	}

	duration, err := time.ParseDuration(subs[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("%s", err)
	}

	return time.Now().Add(duration), nil
}

func CalcLastBackup(interval int, t time.Time) time.Time {
	_, remainder := DivMod(t.Hour(), interval)
	lastBackup := time.Date(t.Year(), t.Month(), t.Day(), t.Hour()-remainder, 0, 0, 0, t.Location())
	return lastBackup
}

func TimeMeasure(desc string) func() {
	start := time.Now()
	return func() {
		log.Printf("%v took %v", desc, time.Since(start))
	}
}

func TimeMeasureRate(desc string, totalBytes int64) func() {
	start := time.Now()
	return func() {
		elapsed := time.Since(start).Seconds()
		rate := math.Ceil(float64(totalBytes) / (elapsed * 1024 * 1024))
		log.Printf("%v took %v seconds, rate %.1f Mib/s", desc, elapsed, rate)
	}
}
