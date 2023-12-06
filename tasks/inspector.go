package tasks

import (
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/utils"
	"github.com/gorilla/websocket"
	"github.com/hibiken/asynq"
	"net/http"
	"time"
)

var inspector = asynq.NewInspector(asynq.RedisClientOpt{Addr: config.RedisAddr})
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func HandleGetTask(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query()["task_id"][0]
	if taskID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "taskID missing or invalid"})
		return
	}

	info, err := inspector.GetTaskInfo(config.TaskQueue, taskID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]any{"msg": "not found"})
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]any{"msg": map[string]any{"info": info.State.String(), "result": string(info.Result)}})
}

func HandleMonitorTaskWS(w http.ResponseWriter, r *http.Request) {
	taskID := r.URL.Query()["task_id"][0]
	if taskID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "taskID missing or invalid"})
		return
	}

	var info *asynq.TaskInfo

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"msg": "error creating a websocket"})
		return
	}
	defer ws.Close()

	for {
		info, err = inspector.GetTaskInfo(config.TaskQueue, taskID)
		if err != nil {
			utils.WriteJSON(w, http.StatusNotFound, map[string]interface{}{"msg": "not found"})
			return
		}

		switch info.State.String() {
		case "completed":
			ws.WriteJSON(map[string]any{"msg": "completed"})
			ws.Close()
		case "pending":
			ws.WriteJSON(map[string]any{"msg": "pending"})
		case "active":
			if err = ws.WriteMessage(1, info.Result); err != nil {
				//if err = ws.WriteMessage(1, info.Result); err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, map[string]interface{}{"msg": err})
				return
			}
		}

		time.Sleep(time.Second * 3)
	}
}
