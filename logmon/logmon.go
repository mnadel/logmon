package logmon

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/kalafut/imohash"
)

type LogMonitor struct {
	config *Configuration
	db     *Database
}

type LogFile struct {
	Filename string
	file     *os.File
	offset   uint64
}

func NewLogMonitor(config *Configuration) *LogMonitor {
	return &LogMonitor{
		config: config,
		db:     NewDatabase(config),
	}
}

func (lm *LogMonitor) Logs() []*LogFile {
	logfiles := make([]*LogFile, 0)

	for _, glob := range lm.config.Logs {
		logs, err := filepath.Glob(glob)
		if err != nil {
			log.Println("error globbing", glob, err.Error())
		} else {
			for _, logpath := range logs {
				file, err := os.Open(logpath)

				if err != nil {
					log.Println("error opening log", logpath, err.Error())
					continue
				}

				prev, err := lm.db.getHash(logpath)
				if err != nil {
					log.Println("error getting hash", logpath, err.Error())
					continue
				}

				fileBytes, err := ioutil.ReadAll(file)
				if err != nil {
					log.Println("error reading", logpath)
					continue
				}

				curr := imohash.Sum(fileBytes)

				if bytes.Compare(prev, curr[:]) != 0 {
					offset, err := lm.db.getOffset(logpath)
					if err != nil {
						log.Println("error getting offset", logpath, err.Error())
						continue
					}

					log.Println(logpath, "has changed")

					logfiles = append(logfiles, &LogFile{
						Filename: logpath,
						file:     file,
						offset:   offset,
					})
				}
			}
		}
	}

	return logfiles
}

func NewLogFile(path string, offset uint64) (*LogFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	return &LogFile{
		file:   file,
		offset: offset,
	}, nil
}

func (lm *LogMonitor) Close() {
	lm.db.Close()
}
