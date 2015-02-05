package main

import (
	"fmt"

	"github.com/ActiveState/tail"
)

func Tail(filename string) {
	seek := &tail.SeekInfo{0, 0}

	config := tail.Config{
		Location: seek,
		Follow:   true,
	}

	t, err := tail.TailFile(filepath, config)

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
