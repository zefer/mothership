package main

import (
	"time"

	"github.com/golang/glog"
	"github.com/zefer/gompd/mpd"
)

const retryDur time.Duration = time.Second * 3

type watchConn struct {
	watcher *mpd.Watcher
	addr    string
}

func newWatchConn(addr string) (*watchConn, error) {
	conn := &watchConn{addr: addr}
	conn.retryConnect()
	go conn.errorLoop()
	go conn.eventLoop()
	return conn, nil
}

func (w *watchConn) connect() error {
	watcher, err := mpd.NewWatcher("tcp", w.addr, "")
	if err != nil {
		return err
	}
	w.watcher = watcher
	return nil
}

func (w *watchConn) retryConnect() {
	for {
		err := w.connect()
		if err == nil {
			glog.Infof("Watcher: connected to %s", w.addr)
			return
		}
		glog.Errorf("Watcher: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
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
