package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"

	"github.com/golang/glog"
	"github.com/zefer/mpd-web/mpd"
)

type PlayListEntry struct {
	Pos  int    `json:"pos"`
	Name string `json:"name"`
}

func PlayListHandler(client *mpd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			playListList(client, w, r)
			return
		} else if r.Method == "POST" {
			playListUpdate(client, w, r)
			return
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

func playListList(client *mpd.Client, w http.ResponseWriter, r *http.Request) {
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

func playListUpdate(client *mpd.Client, w http.ResponseWriter, r *http.Request) {
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