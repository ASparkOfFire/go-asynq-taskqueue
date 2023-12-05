package tasks

import (
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/utils"
	"github.com/hibiken/asynq"
	"net/http"
)

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	i := asynq.NewInspector(asynq.RedisClientOpt{Addr: config.RedisAddr})
	taskID := r.URL.Query()["task_id"][0]

	info, err := i.GetTaskInfo(config.TaskQueue, taskID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]any{"msg": "not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"msg": map[string]any{"info": info.State.String(), "result": string(info.Result)}})
}
