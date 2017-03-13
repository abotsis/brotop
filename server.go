package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/elazarl/go-bindata-assetfs"
	"github.com/rsms/gotalk"
)

var (
	Clients = make(map[*gotalk.Sock]int)
	socksmu sync.RWMutex
)

func onAccept(s *gotalk.Sock) {
	// Keep track of connected sockets
	socksmu.Lock()
	defer socksmu.Unlock()
	Clients[s] = 1

	s.CloseHandler = func(s *gotalk.Sock, _ int) {
		socksmu.Lock()
		defer socksmu.Unlock()
		delete(Clients, s)
	}
}

func Broadcast(name string, in interface{}) {
	socksmu.RLock()
	defer socksmu.RUnlock()

	for s, _ := range Clients {
		s.Notify(name, in)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	body, err := Asset("web/index.html")

	if err != nil {
	}

	t, _ := template.New("index").Parse(fmt.Sprintf("{{define 'Version'}}%s{{end}}", string(body)))

	t.ExecuteTemplate(w, "Version", Version)
}

func StartServer() {
	ws := gotalk.WebSocketHandler()
	ws.OnAccept = onAccept

	http.Handle("/gotalk/", ws)

	// for dev
	// http.Handle("/", http.FileServer(http.Dir("./web/")))

	http.Handle("/",
		http.FileServer(
			&assetfs.AssetFS{
				Asset:     Asset,
				AssetDir:  AssetDir,
				AssetInfo: AssetInfo,
				Prefix:    "web",
			},
		),
	)

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		jdata := make(map[string]string)

		jdata["version"] = Version

		js, err := json.Marshal(jdata)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.Write(js)
	})

	var port string = ":8080"

	if len(*ServerPort) > 0 {
		port = fmt.Sprintf(":%s", ServerPort)
	}

	log.Infof("Server listening - http://%s%s", "127.0.0.1", port)

	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Error(err.Error())
	}
}
