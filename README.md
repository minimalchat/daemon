# Let's Chat daemon

[![Build Status](https://travis-ci.org/mihok/letschat-daemon.svg?branch=master)](https://travis-ci.org/mihok/letschat-daemon)
[![Coverage Status](https://coveralls.io/repos/github/mihok/letschat-daemon/badge.svg?branch=master)](https://coveralls.io/github/mihok/letschat-daemon?branch=master)

---

Let's Chat is an open source live chat system providing live one on one messaging to a website visitor and an operator.

Let's Chat is:
-   **minimal**: simple, lightweight, accessible
-   **extensible**: modular, pluggable, hookable, composable

---

Let's Chat daemon is the central server providing API endpoints for operator extensions like Slack, IRC, etc. It also provides the socket.io endpoints that the web clients connect to when on a Let's Chat enabled website.

### Installation

Download the prebuilt binaries available in the [releases]() section or clone the repo and build using Go `>=1.6`.

```
> curl -L https://github.com/mihok/letschat-daemon/releases/download/v1.0.0/letschat.tar.gz
> tar -zxvf letschat.tar.gz
> cd letschat/bin
> letschat -host 0.0.0.0 -port 8080
```

### Usage

```
> letschat-daemon
letschat-daemon runs the socket and API daemon

Find more information at https://github.com/mihok/letschat-daemon

Flags:
  -host address
        The address to which serve socket and API requests on
  -port number
        The port number to use in conjunction with the host address
```
