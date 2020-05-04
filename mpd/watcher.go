package mpd

import (
	"time"

	log "github.com/sirupsen/logrus"
	gompd "github.com/zefer/gompd/mpd"
)

type watcher struct {
	watcher   *gompd.Watcher
	addr      string
	listeners []func(string)
}

func NewWatcher(addr string) *watcher {
	c := &watcher{
		addr:      addr,
		listeners: []func(string){},
	}
	c.connect()
	return c
}

func (w *watcher) OnStateChange(l func(string)) {
	w.listeners = append(w.listeners, l)
}

func (w *watcher) Close() {
	w.watcher.Close()
}

func (w *watcher) connect() {
	for {
		watcher, err := gompd.NewWatcher("tcp", w.addr, "")
		if err == nil {
			w.watcher = watcher
			go w.errorLoop()
			go w.eventLoop()
			log.Infof("MPD watcher: connected to %s", w.addr)
			return
		}
		log.Errorf("MPD watcher: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
}

func (w *watcher) eventLoop() {
	for subsystem := range w.watcher.Event {
		for _, l := range w.listeners {
			l(subsystem)
		}
	}
}

func (w *watcher) errorLoop() {
	for err := range w.watcher.Error {
		log.Errorf("Watcher: %v", err)
		w.Close()
		w.connect()
	}
}
