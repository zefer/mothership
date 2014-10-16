package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/elazarl/go-bindata-assetfs"
	"github.com/golang/glog"
	"github.com/gorilla/mux"

	"github.com/zefer/mpd-web/mpd"
	"github.com/zefer/mpd-web/websocket"
)

var (
	client  *mpd.Client
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	port    = flag.String("port", ":8080", "listen port")
)

func sendStatus(c *websocket.Conn) {
	b, err := mpdStatus()
	if err != nil {
		return
	}
	c.Send(b)
}

func mpdStateChanged(subsystem string) {
	glog.Info("Changed subsystem:", subsystem)
	b, err := mpdStatus()
	if err != nil {
		glog.Errorln(err)
		return
	}
	websocket.Broadcast(b)
}

func main() {
	flag.Parse()
	glog.Infof("Starting API for MPD at %s.", *mpdAddr)

	// Send the browser the MPD state when they first connect.
	websocket.OnConnect(sendStatus)

	// This watcher notifies us when MPD's state changes, without polling.
	watch := mpd.NewWatcher(*mpdAddr)
	defer watch.Close()
	watch.OnStateChange(mpdStateChanged)

	// This client connection provides an API to MPD's commands.
	client = mpd.NewClient(*mpdAddr)
	defer client.Close()

	r := mux.NewRouter()
	r.HandleFunc("/websocket", websocket.Serve)
	r.HandleFunc("/next", NextHandler)
	r.HandleFunc("/previous", PreviousHandler)
	r.HandleFunc("/play", PlayHandler)
	r.HandleFunc("/pause", PauseHandler)
	r.HandleFunc("/randomOn", RandomOnHandler)
	r.HandleFunc("/randomOff", RandomOffHandler)
	r.HandleFunc("/files", FileListHandler)
	r.HandleFunc("/playlist", PlayListHandler)
	r.HandleFunc("/library/updated", LibraryUpdateHandler)

	// The front-end assets are served from a go-bindata file.
	r.PathPrefix("/").Handler(
		http.FileServer(&assetfs.AssetFS{Asset, AssetDir, ""}),
	)
	http.Handle("/", r)
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

func NextHandler(w http.ResponseWriter, r *http.Request) {
	err := client.C.Next()
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PreviousHandler(w http.ResponseWriter, r *http.Request) {
	err := client.C.Previous()
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PlayHandler(w http.ResponseWriter, r *http.Request) {
	err := client.C.Play(-1)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func PauseHandler(w http.ResponseWriter, r *http.Request) {
	err := client.C.Pause(true)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func random(on bool, w http.ResponseWriter) {
	err := client.C.Random(on)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
func RandomOnHandler(w http.ResponseWriter, r *http.Request) {
	random(true, w)
}
func RandomOffHandler(w http.ResponseWriter, r *http.Request) {
	random(false, w)
}

type FileListEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Base string `json:"base"`
}

func FileListHandler(w http.ResponseWriter, r *http.Request) {
	data, err := client.C.ListInfo(r.FormValue("uri"))
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	out := make([]*FileListEntry, len(data))
	for i, item := range data {
		for _, t := range []string{"file", "directory", "playlist"} {
			if p, ok := item[t]; ok {
				out[i] = &FileListEntry{
					Path: p,
					Type: t,
					Base: path.Base(p),
				}
				break
			}
		}
	}
	b, err := json.Marshal(out)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(b))
}

func PlayListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		playListList(w, r)
		return
	} else if r.Method == "POST" {
		playListUpdate(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

type PlayListEntry struct {
	Pos  int    `json:"pos"`
	Name string `json:"name"`
}

func playListList(w http.ResponseWriter, r *http.Request) {
	data, err := client.C.PlaylistInfo(-1, -1)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	out := make([]*PlayListEntry, len(data))
	for i, item := range data {
		var name string
		if artist, ok := item["Artist"]; ok {
			// Artist - Title
			name = fmt.Sprintf("%s - %s", artist, item["Title"])
		} else if n, ok := item["Name"]; ok {
			// Playlist name.
			name = n
		} else {
			// Default to file name.
			name = path.Base(item["file"])
		}
		out[i] = &PlayListEntry{
			Pos:  i + 1,
			Name: name,
		}
	}
	b, err := json.Marshal(out)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(b))
}

func playListUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body.
	decoder := json.NewDecoder(r.Body)
	var params map[string]interface{}
	err := decoder.Decode(&params)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	uri := params["uri"].(string)
	typ := params["type"].(string)
	replace := params["replace"].(bool)
	play := params["play"].(bool)
	pos := 0
	if uri == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Clear the playlist.
	if replace {
		err := client.C.Clear()
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// To play from the start of the new items in the playlist, we need to get the
	// current playlist position.
	if !replace {
		data, err := client.C.Status()
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pos, err = strconv.Atoi(data["playlistlength"])
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		glog.Infof("pos: %d", pos)
	}

	// Add to the playlist.
	if typ == "playlist" {
		err = client.C.PlaylistLoad(uri, -1, -1)
	} else {
		err = client.C.Add(uri)
	}
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Play.
	if play {
		err := client.C.Play(pos)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
}

func LibraryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if r.Method == "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else {
		// Parse the JSON body.
		decoder := json.NewDecoder(r.Body)
		var params map[string]interface{}
		err := decoder.Decode(&params)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		uri := params["uri"].(string)
		if uri == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		_, err = client.C.Update(uri)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	}
}
