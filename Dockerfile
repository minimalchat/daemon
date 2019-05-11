FROM ubuntu:xenial

RUN mkdir -p /daemon
WORKDIR /daemon

RUN apt update
RUN apt install -y golang ca-certificates

COPY dist/daemon /daemon/server

ENTRYPOINT ["/daemon/server", "-host", "0.0.0.0"]
CMD ["-port", "80"]

