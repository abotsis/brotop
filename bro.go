package main

import (
	"bufio"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"strings"

	"github.com/Unknwon/com"
)

type BroHeader struct {
	Name         string
	Separator    string
	SetSeparator string
	EmptyField   string
	UnsetField   string
	Timestamp    string
	Fields       map[string]string
}

var (
	BRO_LOCATIONS = []string{
		"/opt/bro/logs/current",
		"/usr/local/bro/logs/current",
	}
)

func getValue(line string, sep string) []string {
	return strings.Split(strings.Trim(line, "\n"), "\t")[1:]
}

func GetBroHeader(path string) (*BroHeader, error) {
	header := &BroHeader{}

	data, err := os.Open(path)

	defer data.Close()

	if err != nil {
		return header, err
	}

	scanner := bufio.NewReader(data)
	line, err := scanner.ReadString('\n')

	var fields []string
	var types []string

	for err == nil {
		if strings.HasPrefix(line, "#") {
			if strings.HasPrefix(line, "#separator") {
				sep := strings.Trim(strings.Split(line, " ")[1], "\n")
				sepchar, err := hex.DecodeString(sep[2:])

				if err != nil {
					log.Panic(err)
				}

				header.Separator = string(sepchar)

			} else if strings.HasPrefix(line, "#fields") {
				fields = getValue(line, header.Separator)
			} else if strings.HasPrefix(line, "#types") {
				types = getValue(line, header.Separator)
			} else if strings.HasPrefix(line, "#set_separator") {
				header.SetSeparator = getValue(line, header.Separator)[0]
			} else if strings.HasPrefix(line, "#empty_field") {
				header.EmptyField = getValue(line, header.Separator)[0]
			} else if strings.HasPrefix(line, "#unset_field") {
				header.UnsetField = getValue(line, header.Separator)[0]
			} else if strings.HasPrefix(line, "#path") {
				header.Name = getValue(line, header.Separator)[0]
			} else if strings.HasPrefix(line, "#open") {
				header.Timestamp = getValue(line, header.Separator)[0]
			}

			line, err = scanner.ReadString('\n')

		} else {
			break
		}
	}

	if len(fields) <= 0 && len(types) <= 0 {
		return header, errors.New("Not a bro log file.")
	}

	m := make(map[string]string)

	for i, f := range fields {
		m[f] = types[i]
	}

	header.Fields = m

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

func FindBroLogs() ([]*Wrapper, error) {
	var empty []*Wrapper
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

	var valid []*Wrapper

	for _, file := range paths {

		w := NewWrapper(file)

		if w.Init() {
			valid = append(valid, w)
		}
	}

	return valid, nil
}
