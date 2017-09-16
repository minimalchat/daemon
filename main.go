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
	Port     int
	Host     string

	SSLCertFile string
	SSLKeyFile  string
	SSLPort     int

	CORSOrigin string
}

var config configuration
var needHelp bool

func help() {
	fmt.Println("Minimal Chat live chat API/Socket daemon")
	fmt.Println()
	fmt.Println("Find more information at https://github.com/minimalchat/daemon")
	fmt.Println()

	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func init() {
	// Configuration
	flag.StringVar(&config.SSLCertFile, "ssl-cert", "", "SSL Certificate Filepath")
	flag.StringVar(&config.SSLKeyFile, "ssl-key", "", "SSL Key Filepath")
	flag.IntVar(&config.SSLPort, "ssl-port", 443, "Port used to serve SSL HTTPS and websocket traffic on")
	flag.IntVar(&config.Port, "port", 80, "Port used to serve HTTP and websocket traffic on")
	flag.StringVar(&config.Host, "host", "localhost", "IP to serve http and websocket traffic on")
	flag.StringVar(&config.CORSOrigin, "cors-origin", "http://localhost:3000", "Host to allow cross origin resource sharing (CORS)")
	flag.BoolVar(&needHelp, "h", false, "Get help")
}

func main() {
	// Configuration
	flag.Parse()

	if needHelp {
		help()

		return
	}

	// config.Host = fmt.Sprintf("%s:%d", config.IP, config.Port)

	db := new(store.InMemory)

	// Socket.io
	sock, err := socket.Create(db)

	if err != nil {
		log.Fatal(err)
	}

	go sock.Listen()

	// Server
	server := rest.Listen(db)

	// Socket.io handler
	server.Router.HandlerFunc("GET", "/socket.io/", func(resp http.ResponseWriter, req *http.Request) {
		resp.Header().Set("Access-Control-Allow-Origin", config.CORSOrigin)
		resp.Header().Set("Access-Control-Allow-Credentials", "true")
		// resp.Header().Set("Access-Control-Allow-Headers", "X-Socket-Type")

		sock.ServeHTTP(resp, req)
	})

	// Serve SSL/HTTPS if we can
	if config.SSLCertFile != "" && config.SSLKeyFile != "" {
		log.Println(INFO, "server:", fmt.Sprintf("Listening for SSL on %s:%d ...", config.Host, config.SSLPort))
		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Host, config.SSLPort), config.SSLCertFile, config.SSLKeyFile, server.Router)
	}

	log.Println(INFO, "server:", fmt.Sprintf("Listening on %s:%d ...", config.Host, config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), server.Router))
}
