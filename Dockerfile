FROM ubuntu:xenial

RUN mkdir -p /daemon
WORKDIR /daemon

RUN apt update
RUN apt install -y golang

COPY dist/daemon /daemon/server

CMD ["/daemon/server", "--host 0.0.0.0", "--port 80"]

