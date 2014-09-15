package main

import (
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

const pollPeriod = 1 * time.Second

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func writer(ws *websocket.Conn) {
	pollTicker := time.NewTicker(pollPeriod)
	defer func() {
		pollTicker.Stop()
		ws.Close()
	}()
	for {
		select {
		case <-pollTicker.C:
			// Get mpd status.
			b, err := mpdStatus()
			if err != nil {
				glog.Errorln(err)
				return
			}
			// Send mpd status to client.
			if err := ws.WriteMessage(websocket.TextMessage, b); err != nil {
				glog.Errorln(err)
				return
			}
		}
	}
}

func serveWebsocket(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			glog.Error(err)
		}
		return
	}

	go writer(ws)
}
