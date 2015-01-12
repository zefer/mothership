package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sort"

	"github.com/golang/glog"
	"github.com/zefer/mothership/mpd"
)

type FileListEntry struct {
	Path         string `json:"path"`
	Type         string `json:"type"`
	Base         string `json:"base"`
	LastModified string `json:"lastModified"`
}

type ByDate []*FileListEntry
type ByName []*FileListEntry

func (a ByDate) Len() int           { return len(a) }
func (a ByDate) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByDate) Less(i, j int) bool { return a[i].LastModified < a[j].LastModified }

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Path < a[j].Path }

func FileListHandler(c *mpd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		data, err := c.C.ListInfo(r.FormValue("uri"))
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
						Path:         p,
						Type:         t,
						Base:         path.Base(p),
						LastModified: item["last-modified"],
					}
					break
				}
			}
		}

		// Sort the list by the specified field and direction.
		sortFileList(out, r.FormValue("sort"), r.FormValue("direction"))

		b, err := json.Marshal(out)
		w.Header().Add("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	})
}

func sortFileList(f []*FileListEntry, by string, dir string) {
	var t sort.Interface
	// Sort by date or name?
	if by == "date" {
		t = ByDate(f)
	} else {
		t = ByName(f)
	}
	// Sort asc or desc?
	if dir == "desc" {
		sort.Sort(sort.Reverse(t))
	} else {
		sort.Sort(t)
	}
}

func LibraryUpdateHandler(c *mpd.Client) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" && r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
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
		_, err = c.C.Update(uri)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	})
}
