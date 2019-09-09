package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/Unknwon/com"
	_ "github.com/alecthomas/colour"
	"github.com/hpcloud/tail"
	"gopkg.in/alecthomas/kingpin.v1"
)

const (
	Name    = "brotop"
	Version = "0.4.0"
)

var (
	Debug          = kingpin.Flag("debug", "Enable debug mode.").Bool()
	Flush          = kingpin.Flag("flush", "Remove the BroTop cache. (~/.brotop/brotop.db)").Bool()
	DefaultLogPath = kingpin.Flag("path", "Bro log path.").ExistingDir()
	ServerAddr     = kingpin.Flag("listen", "Web server listener address").Short('l').Default("127.0.0.1").String()
	ServerPort     = kingpin.Flag("port", "Web server port.").Short('p').Default("8080").String()
	Quiet          = kingpin.Flag("quiet", "Remove all output logging.").Short('q').Bool()

	OutputChan = make(chan Message)
	DoneChan   = make(chan bool)
)

func init() {
	kingpin.Version(Version)
	kingpin.Parse()

	// log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stderr)

	if *Debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if *Quiet {
		log.SetLevel(log.FatalLevel)
	}

}

func main() {
	log.Infof("Initializing %s Version: %s.", Name, Version)

	home, err := com.HomeDir()

	if err != nil {
		panic(err)
	}

	log.Debug("Looking for BroTop database.")
	brotopPath := path.Join(home, ".brotop")
	brotopDB := path.Join(brotopPath, "brotop.db")

	if *Flush {
		if com.IsExist(brotopDB) {
			err := os.Remove(brotopDB)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

	os.Mkdir(brotopPath, 0777)
	store, err := NewStore(brotopDB, 0600, 1*time.Second)

	if err != nil {
		log.Fatal(err)
	}

	paths, err := FindBroLogs()

	if err != nil {
		log.Fatal(err)
	}

	log.Debug("Hooking to process signals to capture events.")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Info("Starting Webserver.")
	go StartServer()

	log.Info("Opening log files for capture.")
	for _, path := range paths {
		path.Config.Logger = tail.DiscardingLogger
		path.Config.Follow = true
		path.Config.ReOpen = true
		path.Config.Poll = true

		log.Infof("Loading %s.log at %s", path.Name, path.Path)

		var offset int64 = 0

		log.Debug("   * Checking databade for seek information.")

		value, err := store.Get(path.Path)

		if err == nil {
			offset = com.StrTo(value).MustInt64()
			log.Debugf("   * Found: %d.", offset)
		} else {
			log.Debugf("   * Not Found: %d.", offset)
		}

		path.Config.Location = &tail.SeekInfo{offset, os.SEEK_SET}

		log.Debug("   * Now tailing.")
		go path.Capture()
	}

	for {
		select {
		case sig := <-sigChan:
			if sig.String() == "interrupt" {
				log.Debug("Closing done channel.")
				close(DoneChan)
			}
		case msg := <-OutputChan:

			if msg.Error != nil {
				msg.Self.Close()
				log.Fatal(msg.Error)
			}

			json, err := msg.Json()

			if err == nil {
				Broadcast("event", json)

				log.WithFields(log.Fields{
					"type":   msg.Self.Name,
					"path":   msg.Self.Path,
					"length": len(json),
					"offset": msg.Offset,
				}).Debugf("Sending Json Event (%s)", msg.Self.Name)

				store.Set(msg.Self.Path, fmt.Sprintf("%d", msg.Offset))
			}

		case <-DoneChan:
			log.Info("Closing Open Files...")

			// for _, path := range paths {
			// path.Close()
			// }

			log.Info("Cleaning up...")
			store.Close()
			// tail.Cleanup()

			os.Exit(1)
		}
	}
}
