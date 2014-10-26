package mpd

import (
	"time"

	"github.com/golang/glog"
	gompd "github.com/zefer/gompd/mpd"
)

const retryDur time.Duration = time.Second * 3

type Client struct {
	C    *gompd.Client
	addr string
}

func NewClient(addr string) *Client {
	c := &Client{addr: addr}
	c.connect()
	go c.keepAlive()
	return c
}

func (c *Client) Close() {
	c.C.Close()
}

func (c *Client) connect() {
	for {
		client, err := gompd.Dial("tcp", c.addr)
		if err == nil {
			// There will be an old client to close if this was a reconnect.
			if c.C != nil {
				c.C.Close()
			}
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
			// Leave the old connection open until we have a new one because trying to
			// call commands on a closed client will panic.
			c.connect()
		}
		time.Sleep(retryDur)
	}
}
