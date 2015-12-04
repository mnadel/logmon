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

	for _, logfile := range monitor.Logs() {
		log.Println("queueing", logfile.Filename())

		wg.Add(1)
		go func(f *logmon.LogFile) {
			defer f.Close()

			log.Println("processing", f.Filename())

			f.PublishErrors(ch)

			f.Complete()

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

	log.Printf("%v", errors)
}
