package main

import (
	"encoding/json"

	"github.com/zefer/gompd/mpd"
)

type Statuser interface {
	Status() (mpd.Attrs, error)
	CurrentSong() (mpd.Attrs, error)
}

// Returns MPD's status as a JSON byte. Combines data from multiple MPD queries.
func mpdStatusJSON(c Statuser) ([]byte, error) {
	data, err := c.Status()
	if err != nil {
		return nil, err
	}
	song, err := c.CurrentSong()
	if err != nil {
		return nil, err
	}
	for k, v := range song {
		data[k] = v
	}
	b, err := json.Marshal(data)
	return b, err
}
