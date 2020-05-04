package mpd

import (
	"time"

	log "github.com/sirupsen/logrus"
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
			log.Infof("MPD client: connected to %s", c.addr)
			return
		}
		log.Errorf("MPD client: connect failed. Waiting then retrying. %v", err)
		time.Sleep(retryDur)
	}
}

func (c *Client) keepAlive() {
	for {
		err := c.C.Ping()
		if err != nil {
			log.Errorf("MPD client: ping failed, reconnecting")
			// Leave the old connection open until we have a new one because trying to
			// call commands on a closed client will panic.
			c.connect()
		}
		time.Sleep(retryDur)
	}
}

func (c *Client) Pause(pause bool) error {
	return c.C.Pause(pause)
}

func (c *Client) ListInfo(uri string) ([]gompd.Attrs, error) {
	return c.C.ListInfo(uri)
}

func (c *Client) Update(uri string) (jobID int, err error) {
	return c.C.Update(uri)
}

func (c *Client) Next() error {
	return c.C.Next()
}

func (c *Client) Previous() error {
	return c.C.Previous()
}

func (c *Client) Play(pos int) error {
	return c.C.Play(pos)
}

func (c *Client) Status() (gompd.Attrs, error) {
	return c.C.Status()
}

func (c *Client) PlaylistInfo(start, end int) ([]gompd.Attrs, error) {
	return c.C.PlaylistInfo(start, end)
}

func (c *Client) Clear() error {
	return c.C.Clear()
}

func (c *Client) PlaylistLoad(name string, start, end int) error {
	return c.C.PlaylistLoad(name, start, end)
}

func (c *Client) Add(uri string) error {
	return c.C.Add(uri)
}

func (c *Client) Random(random bool) error {
	return c.C.Random(random)
}
