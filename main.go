package main

import (
  "log"
  "fmt"

  "flag"
  // "io"
  "net/http"

  "github.com/googollee/go-socket.io"
  // "github.com/mihok/lets-chat/operator"
  // "github.com/mihok/lets-chat/person"
 )

// Configuration object
type configuration struct {
  Protocol string
  IP string
  Port int
  Host string
}

// Log levels
const (
  DEBUG string = "DEBUG"
  INFO string = "INFO"
  WARNING string = "WARN"
  ERROR string = "ERROR"
  FATAL string = "FATAL"
)


func main() {

  // Configuration
  var config configuration

  flag.IntVar(&config.Port, "port", 8000, "Port used to serve http and websocket traffic on")
  flag.StringVar(&config.IP, "host", "localhost", "IP to serve http and websocket traffic on")

  config.Host = fmt.Sprintf("%s:%d", config.IP, config.Port)


  // Socket.io
  socket, err := socketio.NewServer(nil)

  if err != nil {
      log.Fatal(err)
  }

  // Socket.io - Connection event
  socket.On("connection", func (sock socketio.Socket) {
    log.Println(DEBUG, "socket:", fmt.Sprintf("Incoming connection %s", sock.Id()))

    // Does this user match a previous fingerprint ?
    //  Does user have cookie?
    //  Does user have known IP?

    // If yes, lets create/update the user
    // If no, lets create new user

    // Create new chat, assign user

    // Message event
    sock.On("client:message", func (msg string) {
      log.Println(DEBUG, "client", fmt.Sprintf("%s: %s", sock.Id(), msg))

      // Save chat
    })

    // Disconnection event
    sock.On("disconnection", func () {
      log.Println(DEBUG, "socket:", fmt.Sprintf("%s disconnected", sock.Id()))

      // Save chat
    })
  })

  socket.On("error", func(so socketio.Socket, err error) {
      log.Println(ERROR, "socket:", err)
  })

  // Socket.io handler
  http.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
    resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
    resp.Header().Set("Access-Control-Allow-Credentials", "true")

    socket.ServeHTTP(resp, req)
  })

  // Server
  log.Println(INFO, "server:", fmt.Sprintf("Listening on %s ...", config.Host))
  log.Fatal(http.ListenAndServe(config.Host, nil))
}