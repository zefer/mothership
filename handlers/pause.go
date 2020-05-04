package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Pauser interface {
	Pause(pause bool) error
}

func PauseHandler(c Pauser) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := c.Pause(true)
		if err != nil {
			log.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
