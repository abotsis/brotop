package main

import (
	"fmt"

	_ "github.com/alecthomas/colour"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	Name    = "brotop"
	Version = "0.1.0"
)

var (
	Debug          = kingpin.Flag("debug", "Enable debug mode.").Bool()
	DefaultLogPath = kingpin.Flag("path", "Bro log path.").ExistingDir()
)

func main() {
	kingpin.Version(Version)
	kingpin.Parse()

	paths, err := FindBroLogs()

	if err != nil {
		fmt.Println("^1red Fail!")
	}

	fmt.Println(paths)

}
