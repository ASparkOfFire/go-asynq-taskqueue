package tasks

import (
	"fmt"
	"github.com/asparkoffire/go-asynq-taskqueue/config"
	"github.com/asparkoffire/go-asynq-taskqueue/utils"
	"github.com/gorilla/mux"
	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

func EnqueueTask(task *asynq.Task) (info *asynq.TaskInfo, err error) {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: config.RedisAddr})

	defer func(client *asynq.Client) { _ = client.Close() }(client)

	info, err = client.Enqueue(task, asynq.Retention(time.Minute*30), asynq.MaxRetry(0))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	return info, nil
}

func HandleCreateTask(w http.ResponseWriter, r *http.Request) {
	model := r.FormValue("model")
	TTA, err := strconv.ParseBool(r.FormValue("tta"))
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "malformed inputs"})
		return
	}

	// create a temporary directory to store and process images
	temp, err := os.MkdirTemp(config.TempDir, "tmp")
	if err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, map[string]any{"msg": "internal server error"})
		logrus.Error("error creating a temporary directory")
		return
	}
	img, imgHeader, err := r.FormFile("image")
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "bad request"})
		logrus.Error(err)
		return
	}

	if imgHeader.Size > 4194304 {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "image size should not exceed 4M"})
		return
	}
	if !utils.CheckAllowedImageExtension(filepath.Ext(imgHeader.Filename)) {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "invalid file extension", "supportedExtensions": []string{".jpg", ".jpeg", ".png"}})
		return
	}
	if model == "" {
		model = "remacri"
	} else if !utils.CheckModel(model) {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "invalid model string", "supportedModels": []string{
			"4x_NMKD-Superscale-SP_178000_G",
			"realesrgan-x4plus-anime",
			"realesrgan-x4plus",
			"RealESRGAN_General_x4_v3",
			"remacri",
			"ultramix_balanced",
			"ultrasharp",
		}})
		return
	}

	defer img.Close()

	imgBin := make([]byte, imgHeader.Size)
	img.Read(imgBin)

	fileName := path.Join(temp, "image"+filepath.Ext(imgHeader.Filename))
	fmt.Println(fileName)
	os.WriteFile(fileName, imgBin, 0644)

	task, err := NewUpscaleTask(fileName, model, TTA)
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
	return
}

func HandleDownloadImage(w http.ResponseWriter, r *http.Request) {
	taskID := mux.Vars(r)["task_id"]
	if taskID == "" {
		utils.WriteJSON(w, http.StatusBadRequest, map[string]any{"msg": "taskID missing or invalid"})
		return
	}

	var info *asynq.TaskInfo

	info, err := inspector.GetTaskInfo(config.TaskQueue, taskID)
	if err != nil {
		utils.WriteJSON(w, http.StatusNotFound, map[string]interface{}{"msg": "not found"})
		return
	}
	if info.State.String() != "completed" {
		utils.WriteJSON(w, http.StatusNotFound, map[string]interface{}{"msg": "not found"})
		return
	}

	outputPath := string(info.Result)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(outputPath)))
	w.Header().Set("Content-Type", utils.GetContentType(outputPath))
	http.ServeFile(w, r, outputPath)

}
