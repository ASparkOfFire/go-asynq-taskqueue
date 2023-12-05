package worker

import (
	"github.com/asparkoffire/go-asynq-taskqueue/tasks"
	"github.com/hibiken/asynq"
	"log"
)

func Worker() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 1},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSleeperTask, tasks.HandleSleeperTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
