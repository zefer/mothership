package handlers

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Player interface {
	Play(pos int) error
}

func PlayHandler(c Player) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := c.Play(-1)
		if err != nil {
			log.Errorln(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}
