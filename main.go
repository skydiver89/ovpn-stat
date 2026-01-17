package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

var VERSION = ""
var GITREV = ""
var BUILDTIME = ""

var configFile = "config.yml"
var config Config
var port string = "8888"

func parseFlags() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-c configFile] [-p port] [-h] [-v]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	var showHelp bool
	flag.BoolVar(&showHelp, "h", false, "Show this help")
	var showVersion bool
	flag.BoolVar(&showVersion, "v", false, "Show version information")
	flag.StringVar(&configFile, "c", "config.yml", "Path to config file")
	flag.StringVar(&port, "p", "8888", "Port to serve")
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
	//router.LoadHTMLGlob("templates/*/*.html")
	router.Static("/images", "./images")
	router.StaticFile("/favicon.ico", "./images/favicon.ico")
	router.GET("/", mainPage)
	router.Run("0.0.0.0:" + port)
}

func mainPage(c *gin.Context) {
	content, err := os.ReadFile(config.Server.Log)
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	scanner := bufio.NewScanner(strings.NewReader(string(content)))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Peer Connection Initiated") {
			fmt.Println(line)
		}
		if strings.Contains(line, "SIGUSR1") {
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln("Error:", err)
	}
	/*
		c.HTML(http.StatusOK, "main.html", gin.H{
			"currentyear": curYear,
			"srvversion":  GITREV,
		})
	*/
}
