package main

import (
	"fmt"
	"time"

	"github.com/ActiveState/tail"
)

type Wrapper struct {
	Name      string
	Fields    []string
	Path      string
	CreatedAt time.Time
	Types     []string
	Follow    bool
	Seek      tail.SeekInfo
	Config    tail.Config
}

func NewWrapper() *Wrapper {

	self := &Wrapper{}

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

func (self *Wrapper) Init() {

}

func (self *Wrapper) Update() {

}

func (self *Wrapper) Close() {

}
