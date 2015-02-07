package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/ActiveState/tail"
	"github.com/Unknwon/com"
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
	ServerPort     = kingpin.Flag("port", "Web server port.").String()

	OutputChan = make(chan Message)
	DoneChan   = make(chan bool)
)

func main() {
	kingpin.Version(Version)
	kingpin.Parse()

	home, err := com.HomeDir()

	if err != nil {
		panic(err)
	}

	brotopPath := path.Join(home, ".brotop")
	os.Mkdir(brotopPath, 0777)
	store, err := NewStore(path.Join(brotopPath, "brotop.db"), 0600, 1*time.Second)

	if err != nil {
		panic(err)
	}

	paths, err := FindBroLogs()

	if err != nil {
		fmt.Println("^1red Fail!")
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	go StartServer()

	for _, path := range paths {
		path.Config.Follow = true
		path.Config.ReOpen = true
		path.Config.Poll = true

		var offset int64 = 0

		value, err := store.Get(path.Path)

		if err == nil {
			offset = com.StrTo(value).MustInt64()
		}

		path.Config.Location = &tail.SeekInfo{offset, os.SEEK_SET}

		go path.Capture()
	}

	for {
		select {
		case sig := <-sigChan:
			if sig.String() == "interrupt" {
				close(DoneChan)
			}
		case msg := <-OutputChan:

			if msg.Error != nil {
				msg.Self.Close()
				panic(msg.Error)
			}

			// fmt.Printf("%s :: %s\n\n", msg.Self.Name, msg.Data)
			// fmt.Println(msg.Json())
			json, err := msg.Json()

			if err != nil {
				panic(err)
			}

			Broadcast("event", json)
			fmt.Println(json)

			store.Set(msg.Self.Path, fmt.Sprintf("%d", msg.Offset))

		case <-DoneChan:
			fmt.Println("\nClosing Open Files...")

			// for _, path := range paths {
			// path.Close()
			// }

			fmt.Println("Cleaning up...")
			store.Close()
			tail.Cleanup()

			os.Exit(1)
		}
	}
}
