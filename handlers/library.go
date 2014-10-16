package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	"github.com/golang/glog"
	"github.com/zefer/mpd-web/mpd"
)

type FileListEntry struct {
	Path string `json:"path"`
	Type string `json:"type"`
	Base string `json:"base"`
}

func FileListHandler(client *mpd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
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
	})
}

func LibraryUpdateHandler(client *mpd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}
