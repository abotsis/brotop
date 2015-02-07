package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ActiveState/tail"
)

var (
	msgmu sync.RWMutex
)

type Wrapper struct {
	Name   string
	Header *BroHeader
	Path   string
	Follow bool
	Config tail.Config
	Tail   *tail.Tail
}

type Message struct {
	Data   string
	Error  error
	Offset int64
	Self   *Wrapper
}

type JsonLine struct {
	Data      []map[string]string `json:"data"`
	Fields    []string            `json:"fields"`
	Type      string              `json:"type"`
	Path      string              `json:"path"`
	Timestamp time.Time           `json:"timestamp"`
}

func (self *Message) Json() (string, error) {

	msgmu.Lock()
	defer msgmu.Unlock()

	jsonLine := JsonLine{
		Type:      self.Self.Name,
		Path:      self.Self.Path,
		Timestamp: time.Now(),
	}

	sep := self.Self.Header.Separator
	data := strings.Split(self.Data, sep)

	section := make([]map[string]string, len(data))

	for i, value := range self.Self.Header.Fields {
		dmap := make(map[string]string)

		key := self.Self.Header.FieldMap[value]

		dmap["value"] = data[i]
		dmap["type"] = key
		dmap["field"] = value

		section[i] = dmap
	}

	jsonLine.Fields = self.Self.Header.Fields
	jsonLine.Data = section

	j, err := json.Marshal(jsonLine)

	return string(j), err
}

func NewWrapper(path string) *Wrapper {

	self := &Wrapper{
		Path: path,
	}

	return self
}

func (self *Wrapper) Capture() {
	var err error

	self.Tail, err = tail.TailFile(self.Path, self.Config)

	if err != nil {
		fmt.Println("word")
	}

	for line := range self.Tail.Lines {

		offset, err := self.Tail.Tell()
		data := line.Text

		if !strings.HasPrefix(data, "#") {
			message := Message{
				Self:   self,
				Offset: offset,
				Error:  err,
				Data:   data,
			}

			if len(data) > 0 {
				OutputChan <- message
			}
		}

	}

	err = self.Tail.Wait()

	if err != nil {
		message := Message{
			Self:  self,
			Error: err,
		}

		OutputChan <- message
	}
}

func (self *Wrapper) Init() bool {
	header, err := GetBroHeader(self.Path)

	if err != nil {
		return false
	}

	self.Header = header
	self.Name = header.Name

	return true
}

func (self *Wrapper) Update() {

}

func (self *Wrapper) Close() {
	self.Tail.Stop()
}
