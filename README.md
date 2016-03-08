A simple websocket server that proxies a telent connect to/from a telnet
server. Intended to help write webclients.

Requirements
============
 * A working Go environment. See: https://golang.org/doc/install
 * gorilla websockets: `go get github.com/gorilla/websocket`
 * kmud's telnet implementation for go: `go get github.com/Cristofori/kmud/telnet`

Running
=======
`go run muck-websocket.go`

The host and port of the websocket server can be specified with the `addr` flag.
The host and port of the proxied muck can be specified with the `muck` flag.

TODO
====
 * Session management/reconnects
 * Impose real origin checks.
