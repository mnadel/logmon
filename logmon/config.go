package logmon

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
	Recipient  string   `json:"recipient"`
	Sender     string   `json:"sender"`
	Smtp       string   `json:"smtp"`
	Subject    string   `json:"subject"`
	Logs       []string `json:"logs"`
	Db         string   `json:"db"`
	ErrorToken string   `json:"errortoken"`
}

func NewConfiguration(stdin bool, filepath string) *Configuration {
	var config *Configuration

	if stdin {
		config = parseConfig(os.Stdin)
	} else {
		file, err := os.Open(filepath)

		if err != nil {
			log.Fatal("error reading file", filepath, err.Error())
		}

		config = parseConfig(file)
	}

	return config
}

func (config *Configuration) LogMonitor() *LogMonitor {
	return NewLogMonitor(config)
}

func parseConfig(reader io.Reader) *Configuration {
	decoder := json.NewDecoder(reader)

	config := &Configuration{
		ErrorToken: "ERROR",
	}

	if err := decoder.Decode(config); err != nil {
		log.Fatal("error parsing config", err.Error())
	}

	return config
}
