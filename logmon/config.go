package logmon

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Configuration struct {
	AlertConfig map[string]string `json:"alert"`
	Db          string            `json:"db"`
	Logs        []string          `json:"logs"`
	ErrorTokens []string          `json:"toks"`
}

func NewConfiguration(stdin bool, filepath string) *Configuration {
	var config *Configuration

	if stdin {
		config = parseConfig(os.Stdin)
	} else {
		file, err := os.Open(filepath)

		if err != nil {
			log.Fatal("error reading file:", filepath, err.Error())
		}

		config = parseConfig(file)
	}

	return config
}

func (config *Configuration) NewEmailAlerter() Alerter {
	return &EmailAlerter{
		From:    config.AlertConfig["from"],
		To:      config.AlertConfig["to"],
		Smtp:    config.AlertConfig["smtp"],
		Subject: config.AlertConfig["subject"],
	}
}

func (config *Configuration) LogMonitor() *LogMonitor {
	return NewLogMonitor(config)
}

func parseConfig(reader io.Reader) *Configuration {
	decoder := json.NewDecoder(reader)

	config := &Configuration{
		ErrorTokens: []string{"ERROR", "FATAL"},
	}

	if err := decoder.Decode(config); err != nil {
		log.Fatal("error parsing config:", err.Error())
	}

	return config
}
