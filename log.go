package main

import (
	"encoding/json"
	"fmt"
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
	Lock     *sync.RWMutex
	Verbose  bool
}

func initLog(verbose bool) *Logger {
	ret := &Logger{}
	ret.Lock = &sync.RWMutex{}
	ret.Recently = make([]LogSample, 60)
	ret.Verbose = verbose
	return ret
}

func (self *Logger) writeLog(entry LogSample) {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	if self.Verbose {
		entryLn, marshalErr := json.Marshal(entry)
		if nil != marshalErr {
			fmt.Printf("! Can't marshal entry\n")
		} else {
			fmt.Printf("> %s\n", entryLn)
		}
	}
	if self.Offset >= len(self.Recently) {
		self.Offset = 0
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
