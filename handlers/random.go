package handlers

import (
	"net/http"

	"github.com/airbrake/glog"
)

type Randomer interface {
	Random(random bool) error
}

func random(c Randomer, on bool, w http.ResponseWriter) {
	err := c.Random(on)
	if err != nil {
		glog.Errorln(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func RandomOnHandler(c Randomer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		random(c, true, w)
	})
}
func RandomOffHandler(c Randomer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		random(c, false, w)
	})
}
