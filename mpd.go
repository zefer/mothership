package main

import (
	"github.com/golang/glog"
	"github.com/zefer/gompd/mpd"
)

type watchConn struct {
	watcher *mpd.Watcher
}

func newWatchConn(addr string) (*watchConn, error) {
	w, err := mpd.NewWatcher("tcp", addr, "")
	if err != nil {
		return nil, err
	}
	conn := &watchConn{w}
	go conn.errorLoop()
	go conn.eventLoop()

	return conn, nil
}

func (w *watchConn) Close() {
	w.watcher.Close()
}

func (w *watchConn) eventLoop() {
	for subsystem := range w.watcher.Event {
		glog.Info("Changed subsystem:", subsystem)
		broadcastStatus()
	}
}

func (w *watchConn) errorLoop() {
	for err := range w.watcher.Error {
		glog.Errorf("Watcher: %v", err)
	}
}
