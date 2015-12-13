# Mothership

[![Build Status](https://circleci.com/gh/zefer/mothership.svg?&style=shield)](https://circleci.com/gh/zefer/mothership)

A web UI for [MPD](http://www.musicpd.org/) built with Go, AngularJS &
WebSockets.

Designed for people who like to browse their music in its original directory
structure.

MPD state changes are broadcasted to all connected clients via WebSockets, which
keeps all users in sync with what is currently playing.

Builds to a single, self-contained binary making it easy to run on any platform.

## Dev usage (does not package assets in binary)

```
# Build the Angular static html front-end app
(cd frontend && grunt build)
# Run the API & serve the static front-end
go build && mothership -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```

## Build (with assets packaged in the binary)

```
# Build the Angular static html front-end app
(cd frontend && grunt build)
# Compile the assets (dev mode used the go-bindata -debug flag)
go-bindata frontend.go -prefix "frontend/dist/" frontend/dist/...
# Build the binary
go build
# Run it
mothership -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```

To cross-compile for a Raspberry Pi use `GOOS=linux GOARM=7 GOARCH=arm go build`

Note that `GOARM=6` should be used for the Raspberry Pi 1 range, `GOARM=7` for
the Raspberry Pi 2 range (released early 2015).

## Work in progress

![UI](https://dl.dropboxusercontent.com/u/89410/mothership.gif)
