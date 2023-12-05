package tasks

import (
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/utils"
	"github.com/hibiken/asynq"
	"log"
	"net/http"
	"strconv"
)

func EnqueueTask(task *asynq.Task) (info *asynq.TaskInfo, err error) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddr})

	defer func(client *asynq.Client) { _ = client.Close() }(client)

	info, err = client.Enqueue(task)
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	return info, nil
}

func HandleTask(w http.ResponseWriter, r *http.Request) {
	d := r.FormValue("duration")

	duration, _ := strconv.Atoi(d)

	task, err := NewSleeperTask(duration)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := EnqueueTask(task)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{"msg": "error enqueuing the task"})
		log.Fatalln("error enqueuing the task: ", err)
	}

	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	b := map[string]any{"msg": map[string]string{"id": info.ID, "queue": info.Queue}}
	utils.WriteJSON(w, http.StatusOK, b)
}
