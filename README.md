# Mothership

A music player UI for [MPD](http://www.musicpd.org/), built with Go, Svelte and WebSockets. Builds to a single, self-contained binary.

## Prerequisites

- [Go](https://golang.org/) 1.22+
- [Node.js](https://nodejs.org/) 20+

## Quick start

```
make install
MPD_ADDR=your-mpd-host:6600 make dev
```

Then open http://localhost:5173.

## Build

```
make build
```

For Raspberry Pi:

```
make build-pi
```

## Run

```
MPD_ADDR=music:6600 ./mothership
```

Or with flags: `./mothership -mpdaddr music:6600 -port :8080`

## Deploy

Transfer the binary and run it. See [bin/deploy](bin/deploy) for a complete example, and [server configuration](https://github.com/zefer/ansible/tree/master/roles/mothership) using Ansible.

## License

MIT. See [LICENSE](LICENSE).
