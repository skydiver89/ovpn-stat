package main

import (
	"flag"
	"fmt"
	"log"
	"os"
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
}
