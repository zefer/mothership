# mpd-web

Simple web UI for mpd.

## Usage

```
# Build the Angular static html front-end app
(cd frontend && grunt)
# Run the API & serve the static front-end
go run main.go -logtostderr=true -mpdaddr=192.168.33.20:6600 -port :8080
# open the app in your browser
open localhost:8080
```
