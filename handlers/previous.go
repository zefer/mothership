package handlers

import (
	"net/http"

	"github.com/airbrake/glog"
)

type Previouser interface {
	Previous() error
}

func PreviousHandler(c Previouser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := c.Previous()
		if err != nil {
			glog.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
