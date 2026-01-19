package main

import (
	"bufio"
	"log"
	"os"
	"strings"
	"time"
)

type timeInterval struct {
	start time.Time
	end   time.Time
}

func readStat() map[string][]timeInterval {
	content, err := os.ReadFile(config.Server.Log)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	stat := make(map[string][]timeInterval)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Peer Connection Initiated") {
			startTime, user, err := parseLogLine(line)
			if err != nil {
				continue
			}
			stat[user] = append(stat[user], timeInterval{startTime, time.Time{}})
		}
		if strings.Contains(line, "SIGUSR1") {
			endTime, user, err := parseLogLine(line)
			if err != nil {
				continue
			}
			_, hasUser := stat[user]
			if !hasUser {
				continue
			}
			stat[user][len(stat[user])-1].end = endTime
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Error:", err)
	}

	return stat
}
