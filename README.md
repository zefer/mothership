# Mothership

[![Build Status](https://circleci.com/gh/zefer/mothership.svg?&style=shield)](https://circleci.com/gh/zefer/mothership)

Mothership is a music player interface for [MPD][mpd],
optimised for browsing your music collection in its original directory
structure.

Mothership is built with Go, AngularJS & WebSockets providing a snappy,
real-time user experience. All connected clients keep the UI in sync with the
player state.

Mothership is cross-platform & extremely portable, building to a single,
self-contained binary with no external dependencies other than an [MPD][mpd]
server to point it to.

![screenshots](https://user-images.githubusercontent.com/101193/28209521-b7a2e07a-688a-11e7-8a47-ecac04bb9844.png)

## Build

To build the Mothership binary, install the [development
prerequisites](#development-prerequisites), then:

```
(cd frontend && grunt build)
go-bindata frontend.go -prefix "frontend/dist/" frontend/dist/...
go build
```

Cross-compilation is achieved by modifying the `go build` command. For example:

* Build for a Raspberry Pi 2: `GOOS=linux GOARM=7 GOARCH=arm go build`
* Build for a Raspberry Pi 1: `GOOS=linux GOARM=6 GOARCH=arm go build`
* Build for linux/386: `GOOS=linux GOARCH=386 go build`
* Build for darwin/386: `GOOS=darwin GOARCH=386 go build`
* Build for windows/386: `GOOS=windows GOARCH=386 go build`

## Run

Firstly, [build Mothership](#build), then:

```
mothership -mpdaddr=music:6600 -port :8080
open localhost:8080
```

Modify `-mpdaddr` to point to the host:port (or ip:port) of your running MPD
server.

## Deploy

Deployment is simple, transfer the binary & run it. A complete example is
provided below:

* [Example server configuration](https://github.com/zefer/ansible/tree/master/roles/mothership)
  (using Ansible)
* [Example deploy script](bin/deploy)

## Develop

While developing, the assets are not packaged into a binary, this allows us to
make front-end changes without rebuilding the back-end.

Install the [development prerequisites](#development-prerequisites), then:

```
(cd frontend && grunt)
go build && mothership -mpdaddr=music:6600 -port :8080
open localhost:8080
```

`grunt` watches for changes and runs all front-end tests.

`go test ./...` runs all the back-end tests, or run a single package with a
command like `(cd handlers && go test)`.

## Development prerequisites

* [Go][go]
* [Node.js][nodejs] & npm
* `go get github.com/jteeuwen/go-bindata/...`
* `(cd frontend && npm install)`

## Extras

* Add an LCD with [Flashlight][flashlight]
* Build a multi-room audio system with Raspberry Pis, [MPD][mpd] &
  [PulseAudio][pulseaudio]:
  * [Server](https://github.com/zefer/ansible/blob/master/music_server.yml)
  * [Clients](https://github.com/zefer/ansible/blob/master/music_client_pulse.yml)

## License

This project uses the MIT License. See [LICENSE](LICENSE).

[MPD]: http://www.musicpd.org/
[go]: https://golang.org/
[nodejs]: https://nodejs.org/
[pulseaudio]: http://www.freedesktop.org/wiki/Software/PulseAudio/
[flashlight]: https://github.com/zefer/flashlight
