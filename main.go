package main

import (
	"flag"
	"github.com/asparkoffire/go-image-upscaler/config"
	"github.com/asparkoffire/go-image-upscaler/transport"
	"github.com/asparkoffire/go-image-upscaler/worker"
	"log"
)

func main() {
	listenAddr := flag.String("listenAddr", config.HTTPAddr, "Listen address for HTTP transport.")
	flag.Parse()

	go func() {
		log.Fatal(transport.StartServer(listenAddr))
	}()
	worker.Worker()
}
