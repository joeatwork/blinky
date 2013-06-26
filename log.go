package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type LogSample struct {
	Red         float64   `json:"Red"`
	RedScaled   float64   `json:"RedScaled"`
	Green       float64   `json:"Green"`
	GreenScaled float64   `json:"GreenScaled"`
	Blue        float64   `json:"Blue"`
	BlueScaled  float64   `json:"BlueScaled"`
	Color       uint32    `json:"Color"`
	When        time.Time `json:"When"`
}

type Logger struct {
	Recently []LogSample
	Offset   int
	LogFile  string
	Lock     *sync.RWMutex
}

func initLog(logName string) *Logger {
	ret := &Logger{}
	ret.Lock = &sync.RWMutex{}
	ret.Recently = make([]LogSample, 60)
	ret.LogFile = logName
	return ret
}

func (self *Logger) writeLog(entry LogSample) {
	self.Lock.Lock()
	defer self.Lock.Unlock()
	if self.Offset >= len(self.Recently) {
		self.Offset = 0
		file, err := os.OpenFile(self.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm&0766)
		if err != nil {
			fmt.Printf("Can't write to log file %s (%s)\n",
				self.LogFile,
				err.Error())
		} else {
			defer file.Close()
			entryLn, marshalErr := json.Marshal(entry)
			if marshalErr != nil {
				fmt.Printf("Can't marshal entry for log: %s\n", err.Error())
			}
			fmt.Printf("WRITING %s\n", entryLn)
			_, werr := file.Write(entryLn)
			if werr != nil {
				fmt.Printf("Can't write entry to log: %s\n", werr.Error())
			}
			_, nerr := file.Write([]byte("\n"))
			if nerr != nil {
				fmt.Printf("Can't write entry terminator: %s\n", nerr.Error())
			}
		}
	}
	self.Recently[self.Offset] = entry
	self.Offset = self.Offset + 1
}

func (self *Logger) recent() (ret []LogSample) {
	self.Lock.RLock()
	defer self.Lock.RUnlock()
	ret = make([]LogSample, len(self.Recently))
	copy(ret, self.Recently)
	return
}

func (self *Logger) ancient() (ret []LogSample) {
	self.Lock.RLock()
	defer self.Lock.RUnlock()

	file, err := os.Open(self.LogFile)
	if err != nil {
		fmt.Printf("Can't read log file %s (%s)\n",
			self.LogFile,
			err)
		return
	}
	defer file.Close()
	latest := make([]LogSample, len(self.Recently))
	scanner := bufio.NewScanner(file)
	element := -1
	for scanner.Scan() {
		element = element + 1
		if element >= len(latest) {
			element = 0
		}

		line := scanner.Bytes()
		var sample LogSample
		jsonerr := json.Unmarshal(line, &sample)
		if jsonerr != nil {
			fmt.Printf("Log line is apparently not JSON%s\n%s\n", jsonerr, line)
		} else {
			latest[element] = sample
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Can't read contents of %s\n", self.LogFile)
		return
	}

	ret = latest
	return
}
