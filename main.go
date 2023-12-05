package main

import (
	"github.com/asparkoffire/go-asynq-taskqueue/transport"
	"github.com/asparkoffire/go-asynq-taskqueue/worker"
	"log"
)

func main() {
	go func() {
		log.Fatal(transport.StartServer())
	}()
	worker.Worker()
}
