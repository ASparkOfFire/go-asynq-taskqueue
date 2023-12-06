package worker

import (
	"github.com/asparkoffire/go-image-upscaler/tasks"
	"github.com/hibiken/asynq"
	"log"
)

func Worker() {
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{Concurrency: 1},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(tasks.TypeSleeperTask, tasks.HandleUpscaleTask)

	if err := srv.Run(mux); err != nil {
		log.Fatal(err)
	}
}
