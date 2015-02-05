package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Unknwon/com"
)

type BroHeader struct {
	Name         string
	Separator    string
	SetSeparator string
	EmptyField   string
	UnsetField   string
	CratedAt     time.Time
	Fields       map[string]string
}

var (
	BRO_LOCATIONS = []string{"/opt/bro/logs/current", "/usr/local/bro/logs/current"}
)

func GetBroHeader(path string) (*BroHeader, error) {
	header := &BroHeader{}

	data, err := os.Open(path)

	defer data.Close()

	if err != nil {
		return header, err
	}

	scanner := bufio.NewReader(data)
	line, err := scanner.ReadString('\n')

	for err == nil {

		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#separator") {
				sep := strings.Split(line, " ")[1]
				fmt.Println("SEP::", sep)
				sepchar, err := hex.DecodeString(sep[2:])
				if err != nil {
					log.Panic(err)
				}
				header.Separator = string(sepchar)
			} else if strings.HasPrefix(line, "#fields") {
				fields := strings.Split(line, "\t")[1:]
				for idx, typ := range fields {
					fmt.Println(idx, typ)
				}
			} else if strings.HasPrefix(line, "#types") {
				types := strings.Split(line, "\t")[1:]
				for idx, typ := range types {
					fmt.Println(idx, typ)
				}
			}

			line, err = scanner.ReadString('\n')

		} else {
			break
		}

	}

	if err != io.EOF {
		return header, err
	}

	return header, nil
}

func CheckDefaultLocations() (bool, string) {
	for _, p := range BRO_LOCATIONS {
		if com.IsExist(p) {
			return true, p
		}
	}

	return false, ""
}

func FindBroLogs() ([]interface{}, error) {
	var empty []interface{}
	var paths []string
	var err error

	if len(*DefaultLogPath) > 0 {
		paths, err = com.GetFileListBySuffix(*DefaultLogPath, ".log")
	} else {
		has, path := CheckDefaultLocations()

		if !has {
			return empty, errors.New("Error: Bro log path not found. Please use the --path switch.")
		}

		paths, err = com.GetFileListBySuffix(path, ".log")
	}

	if err != nil {
		return empty, err
	}

	var valid []interface{}

	for _, file := range paths {

		w := NewWrapper(file)

		if w.Init() {
			valid = append(valid, w)
		}
	}

	return valid, nil
}
