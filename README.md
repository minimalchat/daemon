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

### Installation

Download the prebuilt binaries available in the [releases]() section or clone the repo and build using Go `>=1.6`.

```
> curl -L https://github.com/minimalchat/daemon/releases/download/v1.0.0/mnml.tar.gz
> tar -zxvf mnml.tar.gz
> cd mnml/bin
> mnml -host 0.0.0.0 -port 8080
```

### Usage

```
> daemon
mnml-daemon runs the socket and API daemon

Find more information at https://github.com/minimalchat/mnml-daemon

Flags:
  -host address
        The address to which serve socket and API requests on
  -port number
        The port number to use in conjunction with the host address
```
