package main

import (
	"flag"
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/transport"
	"github.com/asparkoffire/go-asynq-taskqueue/worker"
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
