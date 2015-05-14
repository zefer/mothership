package handlers

import (
	"net/http"

	"gopkg.in/airbrake/glog.v1"
)

type Pauser interface {
	Pause(pause bool) error
}

func PauseHandler(c Pauser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := c.Pause(true)
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
