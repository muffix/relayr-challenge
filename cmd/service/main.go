package main

import (
	"flag"

	"github.com/muffix/go-microservice-template/internal/httpapi"
)

var (
	defaultPort = 8080

	servicePort int
)

func processCommandlineArgs() {
	flag.IntVar(&servicePort, "p", defaultPort, "Port to listen on to serve HTTP requests")
	flag.Parse()
}

func main() {
	processCommandlineArgs()
	httpapi.NewService(servicePort).Start()
}
