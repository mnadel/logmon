package logmon

import (
	"bytes"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kalafut/imohash"
)

type LogMonitor struct {
	config *Configuration
	db     *Database
}

func NewLogMonitor(config *Configuration) *LogMonitor {
	return &LogMonitor{
		config: config,
		db:     NewDatabase(config),
	}
}

func (lm *LogMonitor) IsError(text string) bool {
	for _, tok := range lm.config.ErrorTokens {
		if strings.Index(text, tok) >= 0 {
			return true
		}
	}
	return false
}

func (lm *LogMonitor) Logs() []*LogFile {
	logfiles := make([]*LogFile, 0)

	for _, glob := range lm.config.Logs {
		logs, err := filepath.Glob(glob)
		if err != nil {
			log.Println("error globbing:", glob, err.Error())
		} else {
			for _, logpath := range logs {
				prev, err := lm.db.getHash(logpath)
				if err != nil {
					log.Println("error getting hash:", logpath, err.Error())
					continue
				}

				curr, err := imohash.SumFile(logpath)
				if err != nil {
					log.Println("error summing file:", logpath, err.Error())
					continue
				}

				if bytes.Compare(prev, curr[:]) != 0 {
					offset, err := lm.db.getOffset(logpath)
					if err != nil {
						log.Println("error getting offset:", logpath, err.Error())
						continue
					}

					log.Println("detected changes:", offset, logpath)

					file, err := os.Open(logpath)
					if err != nil {
						log.Println("error opening log:", logpath, err.Error())
						continue
					}

					if _, err := file.Seek(int64(offset), os.SEEK_SET); err != nil {
						log.Println("cannot seek:", logpath, err.Error())
						continue
					}

					logfiles = append(logfiles, &LogFile{
						file:    file,
						monitor: lm,
						hash:    curr[:],
					})
				} else {
					log.Println("no changes:", logpath)
				}
			}
		}
	}

	return logfiles
}

func (lm *LogMonitor) Close() {
	lm.db.Close()
}

func (lm *LogMonitor) updateDb(filepath string, offset uint64, hash []byte) error {
	return lm.db.updateFile(filepath, offset, hash)
}
