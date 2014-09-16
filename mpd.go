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

type clientConn struct {
	c    *mpd.Client
	addr string
}

func newWatchConn(addr string) *watchConn {
	c := &watchConn{addr: addr}
	c.connect()
	return c
}

func newClientConn(addr string) *clientConn {
	c := &clientConn{addr: addr}
	c.connect()
	return c
}

func (w *watchConn) connect() {
	for {
		watcher, err := mpd.NewWatcher("tcp", w.addr, "")
		if err == nil {
			w.watcher = watcher
			go w.errorLoop()
			go w.eventLoop()
			glog.Infof("MPD watcher: connected to %s", w.addr)
			return
		}
		glog.Errorf("MPD watcher: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
}

func (c *clientConn) connect() {
	for {
		client, err := mpd.Dial("tcp", c.addr)
		if err == nil {
			c.c = client
			glog.Infof("MPD client: connected to %s", c.addr)
			return
		}
		glog.Errorf("MPD client: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
}

func (w *watchConn) Close() {
	w.watcher.Close()
}

func (c *clientConn) Close() {
	c.c.Close()
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
		w.Close()
		w.connect()
	}
}
