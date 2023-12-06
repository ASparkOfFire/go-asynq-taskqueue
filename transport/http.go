package transport

import (
	"github.com/asparkoffire/go-asynq-taskqueue/tasks"
	"github.com/asparkoffire/go-asynq-taskqueue/utils"
	"github.com/gorilla/mux"
	"net/http"
)

func StartServer(listenAddr *string) error {
	router := mux.NewRouter()

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusNotFound, map[string]any{"msg": "not found"})
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utils.WriteJSON(w, http.StatusMethodNotAllowed, map[string]any{"msg": "method not allowed"})
	})

	router.HandleFunc("/upscale", tasks.HandleCreateTask).Methods(http.MethodPost)
	router.HandleFunc("/status/{task_id}", tasks.HandleGetTask).Methods(http.MethodGet)
	router.HandleFunc("/monitor/{task_id}", tasks.HandleMonitorTaskWS).Methods(http.MethodGet)
	router.HandleFunc("/download/{task_id}", tasks.HandleDownloadImage).Methods(http.MethodGet)

	if err := http.ListenAndServe(*listenAddr, router); err != nil {
		return err
	}
	return nil
}
