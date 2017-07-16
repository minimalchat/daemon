package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	// "github.com/julienschmidt/httprouter" // Http router
	// "github.com/googollee/go-socket.io" // Socket

	// "github.com/wayn3h0/go-uuid" // UUID (RFC 4122)

	"github.com/minimalchat/daemon/server/rest"
	"github.com/minimalchat/daemon/server/socket"
	"github.com/minimalchat/daemon/store"
	// "github.com/minimalchat/daemon/operator"
	// "github.com/minimalchat/daemon/client"
	// "github.com/minimalchat/daemon/chat"
)

// Log levels
const (
	DEBUG   string = "DEBUG"
	INFO    string = "INFO"
	WARNING string = "WARN"
	ERROR   string = "ERROR"
	FATAL   string = "FATAL"
)

// Configuration object
type configuration struct {
	Protocol string
	IP       string
	Port     int
	Host     string
}

var config configuration
var needHelp bool

func help() {
	fmt.Println("mnml-daemon live chat API daemon")
	fmt.Println()
	fmt.Println("Find more information at https://github.com/minimalchat/mnml-daemon")
	fmt.Println()

	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func init() {
	// Configuration
	flag.IntVar(&config.Port, "port", 8000, "Port used to serve http and websocket traffic on")
	flag.StringVar(&config.IP, "host", "localhost", "IP to serve http and websocket traffic on")
	flag.BoolVar(&needHelp, "h", false, "Get help")
}

func main() {
	// Configuration
	flag.Parse()

	if needHelp {
		help()
		return
	}

	config.Host = fmt.Sprintf("%s:%d", config.IP, config.Port)

	db := new(store.InMemory)

	// Socket.io
	sock, err := socket.Create()

	if err != nil {
		log.Fatal(err)
	}

	go sock.Listen()

	// Server
	server := rest.Listen(config.IP, config.Port, db)

	// Socket.io handler
	server.Router.HandlerFunc("GET", "/socket.io/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		// resp.Header().Set("Access-Control-Allow-Headers", "X-Socket-Type")

		sock.ServeHTTP(resp, req)
	})

	log.Println(INFO, "server:", fmt.Sprintf("Listening on %s ...", config.Host))

	log.Fatal(http.ListenAndServe(config.Host, server.Router))
}
