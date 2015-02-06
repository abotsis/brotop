package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ActiveState/tail"
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
	Data      map[string]map[string]string `json:"data"`
	Type      string                       `json:"type"`
	Path      string                       `json:"path"`
	Timestamp time.Time                    `json:"timestamp"`
}

func (self *Message) Json() (string, error) {

	jsonLine := JsonLine{
		Type:      self.Self.Name,
		Path:      self.Self.Path,
		Timestamp: time.Now(),
	}

	sep := self.Self.Header.Separator
	data := strings.Split(self.Data, sep)

	section := make(map[string]map[string]string)

	var index int = 0

	for key, value := range self.Self.Header.Fields {
		dmap := make(map[string]string)

		dmap["value"] = data[index]
		dmap["type"] = value

		section[key] = dmap

		index += 1
	}

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

		message := Message{
			Self:   self,
			Offset: offset,
			Error:  err,
			Data:   line.Text,
		}

		OutputChan <- message
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
