package main

import (
	"embed"
	"flag"
	"io/fs"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	airbrake "gopkg.in/gemnasium/logrus-airbrake-hook.v4"

	"github.com/zefer/mothership/handlers"
	"github.com/zefer/mothership/mpd"
	"github.com/zefer/mothership/websocket"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

var (
	client  *mpd.Client
	mpdAddr = flag.String("mpdaddr", "127.0.0.1:6600", "MPD address")
	port    = flag.String("port", ":8080", "listen port")

	abProjectID = flag.Int64("abprojectid", 0, "Airbrake project ID")
	abApiKey    = flag.String("abapikey", "", "Airbrake API key")
	abEnv       = flag.String("abenv", "development", "Airbrake environment name")
)

func main() {
	flag.Parse()

	// Allow env var to override the default MPD address, while still
	// respecting an explicit -mpdaddr flag.
	if env := os.Getenv("MPD_ADDR"); env != "" && !isFlagSet("mpdaddr") {
		*mpdAddr = env
	}
	if env := os.Getenv("PORT"); env != "" && !isFlagSet("port") {
		*port = ":" + env
	}
	log.Infof("Starting API for MPD at %s.", *mpdAddr)

	if *abProjectID > int64(0) && *abApiKey != "" {
		log.AddHook(airbrake.NewHook(*abProjectID, *abApiKey, *abEnv))
	}

	// Send the browser the MPD state when they first connect.
	websocket.OnConnect(func(c *websocket.Conn) {
		b, err := mpdStatusJSON(client.C)
		if err == nil {
			c.Send(b)
		}
	})

	// This watcher notifies us when MPD's state changes, without polling.
	watch := mpd.NewWatcher(*mpdAddr)
	defer watch.Close()
	// When mpd state changes, broadcast it to all websockets.
	watch.OnStateChange(func(s string) {
		log.Info("MPD state change in subsystem: ", s)
		b, err := mpdStatusJSON(client.C)
		if err != nil {
			log.Errorln(err)
			return
		}
		websocket.Broadcast(b)
	})

	// This client connection provides an API to MPD's commands.
	client = mpd.NewClient(*mpdAddr)
	defer client.Close()

	http.HandleFunc("/websocket", websocket.Serve)
	http.Handle("/next", handlers.NextHandler(client))
	http.Handle("/previous", handlers.PreviousHandler(client))
	http.Handle("/play", handlers.PlayHandler(client))
	http.Handle("/pause", handlers.PauseHandler(client))
	http.Handle("/randomOn", handlers.RandomOnHandler(client))
	http.Handle("/randomOff", handlers.RandomOffHandler(client))
	http.Handle("/files", handlers.FileListHandler(client))
	http.Handle("/playlist", handlers.PlayListHandler(client))
	http.Handle("/library/updated", handlers.LibraryUpdateHandler(client))

	// Serve the embedded frontend assets.
	distFS, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to create sub filesystem: %s", err)
	}
	http.Handle("/", http.FileServer(http.FS(distFS)))

	log.Infof("Listening on %s.", *port)
	err = http.ListenAndServe(*port, nil)
	if err != nil {
		log.Errorf("http.ListenAndServe %s failed: %s", *port, err)
		return
	}
}

// isFlagSet returns true if the named flag was explicitly set on the command line.
func isFlagSet(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
