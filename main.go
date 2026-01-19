package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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
	//router.LoadHTMLGlob("templates/*/*.html")
	router.Static("/images", "./images")
	router.StaticFile("/favicon.ico", "./images/favicon.ico")
	router.GET("/", mainPage)
	err = router.Run("0.0.0.0:" + config.Server.Port)
	if err != nil {
		log.Fatalln("Cant run server", err)
	}
}

func mainPage(c *gin.Context) {
	stat := readStat()
	fmt.Println(stat)
	/*
		c.HTML(http.StatusOK, "main.html", gin.H{
			"currentyear": curYear,
			"srvversion":  GITREV,
		})
	*/
}
