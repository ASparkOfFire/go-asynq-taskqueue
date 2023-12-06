package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os/exec"
	"path"
	"path/filepath"
	"sync"
	"time"
)

const (
	TypeSleeperTask = "upscaler:upscale"
)

type UpscaleTaskPayload struct {
	Image string
	Model string
	TTA   bool
}

func NewUpscaleTask(img string, model string, TTA bool) (*asynq.Task, error) {
	payload, err := json.Marshal(UpscaleTaskPayload{Image: img, Model: model, TTA: TTA})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSleeperTask, payload, asynq.Retention(time.Hour*24)), nil
}

func HandleUpscaleTask(ctx context.Context, t *asynq.Task) error {
	var payload UpscaleTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json unmarshal error")
	}

	tta := payload.TTA
	model := payload.Model

	outputPath := path.Join(path.Dir(payload.Image), "image_upsacled"+filepath.Ext(payload.Image))
	cmd := exec.Command("./upscale/upscayl", "-i", payload.Image, "-o", outputPath, "-n", model)

	if tta {
		cmd.Args = append(cmd.Args, "-x")
	}

	fmt.Println(cmd)

	stdout, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		log.Fatal("Error running command")
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		var output []byte
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			//fmt.Printf("[%s] %s\n", "stdout", scanner.Text())
			output = append(append(output, '\n'), scanner.Bytes()...)
			//fmt.Println(string(output))
			_, _ = t.ResultWriter().Write(output)
		}
	}()

	wg.Wait()

	// empty the result to save memory
	_, _ = t.ResultWriter().Write([]byte(outputPath))

	return nil
}
