package logmon

import (
	"bufio"
	"log"
	"os"
)

type LogFile struct {
	file    *os.File
	monitor *LogMonitor
	hash    []byte
}

type LogError struct {
	Text     string
	Filename string
}

func (f *LogFile) Filename() string {
	return f.file.Name()
}

func (f *LogFile) Close() {
	f.file.Close()
}

func (f *LogFile) Complete() {
	log.Println("updating:", f.Filename())

	err := f.monitor.updateDb(f.Filename(), f.Offset(), f.hash)
	if err != nil {
		log.Println("error updating:", f.Filename(), err.Error())
	}

	f.Close()
}

func (f *LogFile) Offset() uint64 {
	pos, err := f.file.Seek(0, os.SEEK_CUR)
	if err != nil {
		log.Println("error getting offset:", f.Filename(), err.Error())
	}

	return uint64(pos)
}

func (f *LogFile) PublishErrors(ch chan<- *LogError) {
	scanner := bufio.NewScanner(f.file)

	log.Println("scanning:", f.Filename())

	for scanner.Scan() {
		line := scanner.Text()

		if f.monitor.IsError(line) {
			errlog := &LogError{
				Text:     line,
				Filename: f.Filename(),
			}

			ch <- errlog
		}
	}
}
