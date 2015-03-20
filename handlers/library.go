package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/golang/glog"
)

type LibraryUpdater interface {
	Update(string) (int, error)
}

func LibraryUpdateHandler(c LibraryUpdater) http.Handler {
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
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if _, ok := params["uri"]; !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		uri := params["uri"].(string)
		if uri == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err = c.Update(uri)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		return
	})
}
