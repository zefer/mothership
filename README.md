# Mothership

A web UI for [MPD](http://www.musicpd.org/) built with Go, AngularJS & WebSockets.

Builds to a single, self-contained binary making it easy to run on any platform.

## Dev usage (does not package assets in binary)

```
# Build the Angular static html front-end app
(cd frontend && grunt)
# Run the API & serve the static front-end
go run *.go -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# Or build the binary & run that
go build && mothership -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```

## Build (with assets packaged in the binary)

```
# Build the Angular static html front-end app
(cd frontend && grunt)
# Compile the assets (dev mode used the go-bindata -debug flag)
go-bindata frontend.go -prefix "frontend/dist/" frontend/dist/...
# Build the binary
go build
# Run it
mothership -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```

To cross-compile for a Raspberry Pi use `GOOS=linux GOARM=6 GOARCH=arm go build`

## Work in progress

![UI](https://dl.dropboxusercontent.com/u/89410/player.gif)
