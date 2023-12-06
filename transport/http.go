package transport

import (
	"github.com/asparkoffire/go-asynq-taskqueue/tasks"
	"net/http"
)

func StartServer(listenAddr *string) error {
	http.HandleFunc("/task", tasks.HandleTask)
	http.HandleFunc("/status", tasks.HandleGetTask)
	http.HandleFunc("/monitor", tasks.HandleMonitorTaskWS)
	http.HandleFunc("/download", tasks.HandleDownloadImage)

	if err := http.ListenAndServe(*listenAddr, nil); err != nil {
		return err
	}
	return nil
}
