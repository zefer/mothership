package mpd

import (
	"time"

	"github.com/golang/glog"
	gompd "github.com/zefer/gompd/mpd"
)

const retryDur time.Duration = time.Second * 3

type watcher struct {
	watcher   *gompd.Watcher
	addr      string
	listeners []func(string)
}

type Client struct {
	C    *gompd.Client
	addr string
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

func NewClient(addr string) *Client {
	c := &Client{addr: addr}
	c.connect()
	go c.keepAlive()
	return c
}

func (w *watcher) connect() {
	for {
		watcher, err := gompd.NewWatcher("tcp", w.addr, "")
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

func (c *Client) connect() {
	for {
		client, err := gompd.Dial("tcp", c.addr)
		if err == nil {
			c.C = client
			glog.Infof("MPD client: connected to %s", c.addr)
			return
		}
		glog.Errorf("MPD client: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
}

func (c *Client) keepAlive() {
	for {
		err := c.C.Ping()
		if err != nil {
			glog.Errorf("MPD client: ping failed, reconnecting")
			c.Close()
			c.connect()
		}
		time.Sleep(retryDur)
	}
}

func (w *watcher) Close() {
	w.watcher.Close()
}

func (c *Client) Close() {
	c.C.Close()
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
		glog.Errorf("Watcher: %v", err)
		w.Close()
		w.connect()
	}
}
