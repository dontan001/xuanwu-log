package util

import (
	"fmt"
	"regexp"
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
