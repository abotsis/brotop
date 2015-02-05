package main

import (
	"fmt"

	"github.com/ActiveState/tail"
)

type Wrapper struct {
	Name   string
	Header *BroHeader
	Path   string
	Follow bool
	Seek   tail.SeekInfo
	Config tail.Config
}

func NewWrapper(path string) *Wrapper {

	self := &Wrapper{
		Path: path,
	}

	return self
}

func (self *Wrapper) Tail() {
	t, err := tail.TailFile(self.Path, self.Config)

	if err != nil {
		fmt.Println("word")
	}

	for line := range t.Lines {
		fmt.Println(line.Text)
		offset, err := t.Tell()
		if err != nil {
			fmt.Println("Error: WTF", err)
		}

		fmt.Println(offset)
	}
}

func (self *Wrapper) Init() bool {
	header, err := GetBroHeader(self.Path)

	if err != nil {
		return false
	}

	self.Header = header

	return true
}

func (self *Wrapper) Update() {

}

func (self *Wrapper) Close() {

}
