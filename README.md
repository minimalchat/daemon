# Minimal Chat daemon

[![GoDoc](https://godoc.org/github.com/minimalchat/daemon?status.svg)](https://godoc.org/github.com/minimalchat/daemon)
[![Build Status](https://travis-ci.org/minimalchat/daemon.svg?branch=master)](https://travis-ci.org/minimalchat/daemon)
[![Coverage Status](https://coveralls.io/repos/github/minimalchat/daemon/badge.svg?branch=master)](https://coveralls.io/github/minimalchat/daemon?branch=master)

---

Minimal Chat is an open source live chat system providing live one on one messaging to a website visitor and an operator.

Minimal Chat is:
-   **minimal**: simple, lightweight, accessible
-   **extensible**: modular, pluggable, hookable, composable 

---

Minimal Chat daemon is the central server providing API endpoints for operator extensions like Slack, IRC, etc. It also provides the socket.io endpoints that the web clients connect to when on a Minimal Chat enabled website.

We're glad you're interested in contributing, feel free to create an [issue](https://github.com/minimalchat/daemon/issues/new) or pick one up but first check out our [contributing doc](https://github.com/minimalchat/daemon/blob/master/CONTRIBUTING.md) and [code of conduct](https://github.com/minimalchat/daemon/blob/master/CODE_OF_CONDUCT.md).


### Installation

Download the prebuilt binaries available in the [releases](https://github.com/minimalchat/daemon/releases) section or clone the repo and build using Go `>=1.6`.

```
> curl -L https://github.com/minimalchat/daemon/releases/download/v0.2.0/daemon-v0.2.0 -o daemon
> chmod +x ./daemon
> ./daemon -host 0.0.0.0 -port 8080
```

### Usage

```
> daemon
Minimal Chat live chat API/Socket daemon

Find more information at https://github.com/minimalchat/daemon

Flags:
  -cors
    	Set if the daemon will handle CORS
  -cors-origin string
    	Host to allow cross origin resource sharing (CORS) (default "http://localhost:3000")
  -h	Get help
  -host string
    	IP to serve http and websocket traffic on (default "localhost")
  -port int
    	Port used to serve HTTP and websocket traffic on (default 8000)
  -ssl-cert string
    	SSL Certificate Filepath
  -ssl-key string
    	SSL Key Filepath
  -ssl-port int
    	Port used to serve SSL HTTPS and websocket traffic on (default 4443)

```
