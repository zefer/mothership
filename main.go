package main

import (
	"encoding/json"
	"flag"
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/glog"

	"github.com/zefer/mothership/handlers"
	"github.com/zefer/mothership/mpd"
	"github.com/zefer/mothership/websocket"
)

var (
	client  *mpd.Client
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	port    = flag.String("port", ":8080", "listen port")
)

func main() {
	flag.Parse()
	glog.Infof("Starting API for MPD at %s.", *mpdAddr)

	// Send the browser the MPD state when they first connect.
	websocket.OnConnect(func(c *websocket.Conn) {
		b, err := mpdStatus()
		if err == nil {
			c.Send(b)
		}
	})

	// This watcher notifies us when MPD's state changes, without polling.
	watch := mpd.NewWatcher(*mpdAddr)
	defer watch.Close()
	// When mpd state changes, broadcast it to all websockets.
	watch.OnStateChange(func(s string) {
		glog.Info("MPD state change in subsystem: ", s)
		b, err := mpdStatus()
		if err != nil {
			glog.Errorln(err)
			return
		}
		websocket.Broadcast(b)
	})

	// This client connection provides an API to MPD's commands.
	client = mpd.NewClient(*mpdAddr)
	defer client.Close()

	http.HandleFunc("/websocket", websocket.Serve)
	http.Handle("/next", handlers.NextHandler(client))
	http.Handle("/previous", handlers.PreviousHandler(client))
	http.Handle("/play", handlers.PlayHandler(client))
	http.Handle("/pause", handlers.PauseHandler(client))
	http.Handle("/randomOn", handlers.RandomOnHandler(client))
	http.Handle("/randomOff", handlers.RandomOffHandler(client))
	http.Handle("/files", handlers.FileListHandler(client))
	http.Handle("/playlist", handlers.PlayListHandler(client))
	http.Handle("/library/updated", handlers.LibraryUpdateHandler(client))

	// The front-end assets are served from a go-bindata file.
	http.Handle("/", http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}))

	glog.Infof("Listening on %s.", *port)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		glog.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

func mpdStatus() ([]byte, error) {
	data, err := client.C.Status()
	if err != nil {
		return nil, err
	}
	song, err := client.C.CurrentSong()
	if err != nil {
		return nil, err
	}
	for k, v := range song {
		data[k] = v
	}
	b, err := json.Marshal(data)
	return b, err
}
