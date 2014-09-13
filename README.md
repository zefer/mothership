# mpd-web

Simple web UI for mpd built with AngularJS & Go.

Compiles to a single self-contained binary for easy deployment.

## Usage

```
# Build the Angular static html front-end app
(cd frontend && grunt)
# Run the API & serve the static front-end
go run *.go -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# Or build the binary & run that
go build && mpd-web -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```

## Cross-compile for Raspberry Pi

`GOOS=linux GOARM=6 GOARCH=arm go build`

## Work in progress

![UI](https://dl.dropboxusercontent.com/u/89410/ui.png)
