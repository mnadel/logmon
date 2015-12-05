package main

import (
	"flag"
	"log"
	"sync"

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
	monitor := config.LogMonitor()
	defer monitor.Close()

	log.Printf("using config: %v", config)

	var wg sync.WaitGroup
	ch := make(chan *logmon.LogError)

	logfiles := monitor.Logs()

	files := make([]*logmon.LogFile, len(logfiles))

	for _, logfile := range logfiles {
		log.Println("queueing:", logfile.Filename())

		wg.Add(1)
		go func(f *logmon.LogFile) {
			log.Println("processing:", f.Filename())
			files = append(files, f)
			f.PublishErrors(ch)
			wg.Done()
		}(logfile)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	errors := make(map[string][]string)

	for err := range ch {
		if errs, ok := errors[err.Filename]; ok {
			errs = append(errs, err.Text)
		} else {
			errors[err.Filename] = []string{err.Text}
		}
	}

	if len(errors) > 0 {
		log.Printf("alerting: %v", errors)

		if err := config.NewEmailAlerter().SendAlert(errors); err != nil {
			log.Fatalln("error alerting:", err.Error())
		}
	}

	for _, f := range files {
		f.Complete()
	}
}
