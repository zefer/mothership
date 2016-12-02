package handlers

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"path"
	"strconv"

	"github.com/zefer/gompd/mpd"
	"gopkg.in/airbrake/glog.v2"
)

type Playlister interface {
	Status() (mpd.Attrs, error)
	PlaylistInfo(int, int) ([]mpd.Attrs, error)
	Clear() error
	PlaylistLoad(string, int, int) error
	Add(string) error
	Play(int) error
}

type PlayListEntry struct {
	Pos  int    `json:"pos"`
	Name string `json:"name"`
}

func PlayListHandler(c Playlister) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			playListList(c, w, r)
			return
		} else if r.Method == "POST" {
			playListUpdate(c, w, r)
			return
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
	})
}

// Helper that returns a start/end range to query from the current playlist.
// We want the full playlist unless it is huge, in which case we want a small
// chunk of it.
func playlistRange(c Playlister) ([2]int, error) {
	// Don't fetch or display more than this many playlist entries.
	max := 500
	// Give MPD a start/end of -1/-1 to fetch the full playlist.
	var rngAll = [2]int{-1, -1}
	// A sensible default for when we can't determine the current playlist size.
	var rngFirst = [2]int{0, max}

	s, err := c.Status()
	if err != nil {
		return rngFirst, err
	}
	if _, ok := s["song"]; !ok {
		// No current song playing, so use the whole (empty) playlist.
		return rngAll, nil
	}

	pos, err := strconv.Atoi(s["song"])
	if err != nil {
		return rngFirst, err
	}
	length, err := strconv.Atoi(s["playlistlength"])
	if err != nil {
		return rngFirst, err
	}

	if length > max {
		// Fetch this chunk of the current playlist. Adjust the starting position to
		// return n items before the current song, for context.
		start := math.Max(float64(0), float64(pos-1))
		return [2]int{int(start), pos + max}, nil
	} else {
		// Fetch all of the current playlist.
		return rngAll, nil
	}
}

func playListList(c Playlister, w http.ResponseWriter, r *http.Request) {
	rng, err := playlistRange(c)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Fetch all, or a slice of the current playlist.
	data, err := c.PlaylistInfo(rng[0], rng[1])
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
		p, err := strconv.Atoi(item["Pos"])
		if err != nil {
			p = 1
		}
		out[i] = &PlayListEntry{
			Pos:  p + 1,
			Name: name,
		}
	}
	b, err := json.Marshal(out)
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprint(w, string(b))
}

func playListUpdate(c Playlister, w http.ResponseWriter, r *http.Request) {
	// Parse the JSON body.
	decoder := json.NewDecoder(r.Body)
	var params map[string]interface{}
	err := decoder.Decode(&params)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check for required fields.
	for _, f := range []string{"uri", "type", "replace", "play"} {
		if _, ok := params[f]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
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
		err := c.Clear()
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// To play from the start of the new items in the playlist, we need to get the
	// current playlist position.
	if !replace {
		data, err := c.Status()
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
		err = c.PlaylistLoad(uri, -1, -1)
	} else {
		err = c.Add(uri)
	}
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Play.
	if play {
		err := c.Play(pos)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
