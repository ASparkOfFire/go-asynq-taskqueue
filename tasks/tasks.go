package tasks

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os/exec"
	"sync"
	"time"
)

const (
	TypeSleeperTask = "sleeper:sleep"
)

type SleeperTaskPayload struct {
	Duration int
}

func NewSleeperTask(duration int) (*asynq.Task, error) {
	payload, err := json.Marshal(SleeperTaskPayload{Duration: duration})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeSleeperTask, payload, asynq.Retention(time.Hour*24)), nil
}

func HandleSleeperTask(_ context.Context, t *asynq.Task) error {
	var payload SleeperTaskPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json unmarshal error")
	}
	fmt.Printf("[*] Sleeping for %v seconds.\n", payload.Duration)

	cmd := exec.Command("./sleeper")

	stdout, err := cmd.StdoutPipe()
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
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Printf("[%s] %s\n", "stdout", scanner.Text())
			log.Fatal(t.ResultWriter().Write([]byte(scanner.Text())))
		}
	}()

	wg.Wait()

	return nil
}
