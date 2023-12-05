package transport

import (
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/tasks"
	"net/http"
)

func StartServer() error {
	http.HandleFunc("/task", tasks.HandleTask)
	http.HandleFunc("/status", tasks.HandleGetTask)
	if err := http.ListenAndServe(config.HTTPAddr, nil); err != nil {
		return err
	}
	return nil
}
