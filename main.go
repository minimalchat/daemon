package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/minimalchat/daemon/server/rest"
	"github.com/minimalchat/daemon/store"
)

// Log levels
const (
	DEBUG   string = "DEBUG"
	INFO    string = "INFO"
	WARNING string = "WARN"
	ERROR   string = "ERROR"
	FATAL   string = "FATAL"
)

var config rest.ServerConfig
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
	flag.IntVar(&config.SSLPort, "ssl-port", 4443, "Port used to serve SSL HTTPS and websocket traffic on")
	flag.IntVar(&config.Port, "port", 8000, "Port used to serve HTTP and websocket traffic on")
	flag.StringVar(&config.Host, "host", "localhost", "IP to serve http and websocket traffic on")
	flag.StringVar(&config.CORSOrigin, "cors-origin", "http://localhost:3000", "Host to allow cross origin resource sharing (CORS)")
	flag.BoolVar(&config.CORSEnabled, "cors", false, "Set if the daemon will handle CORS")
	flag.BoolVar(&needHelp, "h", false, "Get help")
}

func main() {
	// Configuration
	flag.Parse()

	if needHelp {
		help()

		return
	}

	// Create DataStore
	db := new(store.InMemory)

	// Server
	server := rest.Initialize(db, config)

	// Serve SSL/HTTPS if we can
	if config.SSLCertFile != "" && config.SSLKeyFile != "" {
		log.Println(INFO, "server:", fmt.Sprintf("Listening for SSL on %s:%d ...", config.Host, config.SSLPort))
		go http.ListenAndServeTLS(fmt.Sprintf("%s:%d", config.Host, config.SSLPort), config.SSLCertFile, config.SSLKeyFile, server.Router)
	}

	log.Println(INFO, "server:", fmt.Sprintf("Listening on %s:%d ...", config.Host, config.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", config.Host, config.Port), server))
}
