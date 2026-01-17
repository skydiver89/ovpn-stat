package main

import (
	"fmt"
	"regexp"
	"time"
)

func parseLogLine(logLine string) (time.Time, string, error) {
	re := regexp.MustCompile(`^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+(?:[^\s]+\s+)?\[([^\]]+)\]|^(\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\s+([^/\s]+)/`)

	matches := re.FindStringSubmatch(logLine)
	if matches == nil {
		return time.Time{}, "", fmt.Errorf("cant parse line: %s", logLine)
	}

	var timeStr, username string

	if matches[1] != "" {
		timeStr = matches[1]
		username = matches[2]
	} else if matches[3] != "" {
		timeStr = matches[3]
		username = matches[4]
	} else {
		return time.Time{}, "", fmt.Errorf("unknown format")
	}

	// Парсим дату/время
	parsedTime, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Time{}, "", fmt.Errorf("cant parse time: %v", err)
	}

	return parsedTime, username, nil
}
