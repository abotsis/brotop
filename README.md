# brotop - Top for bro logs.

[![Build Status](https://drone.io/github.com/mephux/brotop/status.png)](https://drone.io/github.com/mephux/brotop/latest)

![BroTop](https://raw.githubusercontent.com/mephux/brotop/master/brotop.png)

Brotop lets you stream you bro logs to the browser for easy 
debugging and a real-time glimpse into whats being processed.

# Features

  - 100% dep free binary
  - Auto-detect log locations. If BroTop can't find them use the --path switch.

# Download

  soon.

# Usage

  Just run `brotop` and everything would work. 
  Then open your browser to the port you set. (default port is 8080)

# Building

  Make sure you have go installed.

  - `go get github.com/tools/godep`
  - `go get github.com/jteeuwen/go-bindata/...`
  - `make`

# Package as dep or rpm

  `make package`

