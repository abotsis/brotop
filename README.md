# brotop - Top for bro logs.

[![Build Status](http://104.236.125.70/api/badge/github.com/criticalstack/brotop/status.svg?branch=master)](http://104.236.125.70/github.com/criticalstack/brotop)

<img height="300px" width="" src="https://raw.githubusercontent.com/mephux/brotop/master/brotop.png?token=AABXAe5HY1UJns_gRtUyUvLqkMRYtnRAks5U4Ou_wA%3D%3D">

Brotop lets you stream your bro logs to the browser for easy 
debugging and a real-time glimpse into whats being processed.

# Features

  - 100% dep free binary
  - Auto-detect log locations. If BroTop can't find them use the --path switch.

# Download

  * Linux
    * [brotop linux/amd64](https://github.com/criticalstack/brotop/releases/download/v0.2.1/brotop-linux-amd64.tar.gz)
    * [brotop linux/386](https://github.com/criticalstack/brotop/releases/download/v0.2.1/brotop-linux-386.tar.gz)

  * Macosx
    * [brotop darwin/amd64](https://github.com/criticalstack/brotop/releases/download/v0.2.1/brotop-linux-amd64.tar.gz)
    * [brotop darwin/386](https://github.com/criticalstack/brotop/releases/download/v0.2.1/brotop-linux-386.tar.gz)

  * Freebsd
    * [brotop freebsd/amd64](https://github.com/criticalstack/brotop/releases/download/v0.2.1/brotop-linux-amd64.tar.gz)
    * [brotop freebsd/386](https://github.com/criticalstack/brotop/deleases/download/v0.2.1/brotop-linux-386.tar.gz)

# Usage

  Just run `brotop` and everything would work. 
  Then open your browser to the port you set. (default port is 8080)

  ```
usage: brotop [<flags>]

Flags:
  --help           Show help.
  --debug          Enable debug mode.
  --path=PATH      Bro log path.
  -p, --port=PORT  Web server port.
  -q, --quiet      Remove all output logging.
  --version        Show application version.
  ```

# Building

  Make sure you have go installed.

  - `go get github.com/tools/godep`
  - `go get github.com/jteeuwen/go-bindata/...`
  - `make`

# Package as dep or rpm

  `make package`

