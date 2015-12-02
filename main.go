package main

import (
	"flag"
	"log"

	"./logmon"
)

var stdin = flag.Bool("stdin", false, "Read config from stdin")
var configFile = flag.String("config", "", "Path to config file")

func main() {
	flag.Parse()

	if !*stdin && *configFile == "" {
		log.Fatal("please specify a configuration")
	}

	config := logmon.NewConfiguration(*stdin, *configFile)
	logmon := config.LogMonitor()
	defer logmon.Close()

	log.Printf("using config: %v", config)

	for _, logfile := range logmon.Logs() {
		log.Printf("inspecting: %s", logfile.Filename)
	}
}
