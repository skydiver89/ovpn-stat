package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var VERSION = ""
var GITREV = ""
var BUILDTIME = ""

var configFile = "config.yml"
var config Config

func parseFlags() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-c configFile] [-h] [-v]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "Show this help")
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.StringVar(&configFile, "c", "config.yml", "Path to config file")
	flag.Parse()
	if showHelp {
		flag.Usage()
	}
	if showVersion {
		fmt.Println("Version      : ", VERSION)
		fmt.Println("Git revision : ", GITREV)
		fmt.Println("Build date   : ", BUILDTIME)
		os.Exit(0)
	}
}

func main() {
	parseFlags()
	err := config.load(configFile)
	if err != nil {
		log.Fatalln(err)
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*.html")
	router.Static("/images", "./images")
	router.StaticFile("/favicon.ico", "./images/favicon.ico")
	router.GET("/", mainPage)
	err = router.Run("0.0.0.0:" + config.Server.Port)
	if err != nil {
		log.Fatalln("Cant run server", err)
	}
}

// TemplateData структуры для шаблона
type TemplateData struct {
	Stats map[string]UserSessions
}

type UserSessions struct {
	HasActive bool
	Intervals []SessionInterval
}

type SessionInterval struct {
	IsActive       bool
	StartFormatted string
	EndFormatted   string
	Duration       string
}

func formatDuration(start, end time.Time) string {
	var duration time.Duration
	if end.IsZero() {
		duration = time.Since(start)
	} else {
		duration = end.Sub(start)
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%dч %dм %dс", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%dм %dс", minutes, seconds)
	}
	return fmt.Sprintf("%dс", seconds)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "—"
	}
	return t.Format("02.01.2006 15:04:05")
}

func prepareTemplateData(stat map[string][]timeInterval) TemplateData {
	templateStats := make(map[string]UserSessions)

	for user, intervals := range stat {
		userSessions := UserSessions{
			Intervals: make([]SessionInterval, 0, len(intervals)),
		}

		for _, interval := range intervals {
			isActive := interval.end.IsZero()
			if isActive {
				userSessions.HasActive = true
			}

			sessionInterval := SessionInterval{
				IsActive:       isActive,
				StartFormatted: formatTime(interval.start),
				EndFormatted:   formatTime(interval.end),
				Duration:       formatDuration(interval.start, interval.end),
			}

			userSessions.Intervals = append(userSessions.Intervals, sessionInterval)
		}

		templateStats[user] = userSessions
	}

	return TemplateData{Stats: templateStats}
}

func mainPage(c *gin.Context) {
	stat := readStat()
	templateData := prepareTemplateData(stat)
	c.HTML(http.StatusOK, "index.html", templateData)
}
