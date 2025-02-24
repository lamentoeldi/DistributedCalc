package orchestrator

import (
	"DistributedCalc/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Orchestrator struct {
	Client  *http.Client
	Url     string // protocol://host:port
	Retries int
}

func NewOrchestrator(client *http.Client, url string, retries int) *Orchestrator {
	return &Orchestrator{
		Client:  client,
		Url:     url,
		Retries: retries,
	}
}

func (o *Orchestrator) GetTask(ctx context.Context) (*models.Task, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", o.Url+"/internal/task", nil)
	if err != nil {
		return nil, err
	}

	res, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusNotFound:
			return nil, models.ErrNoTasks
		default:
			return nil, fmt.Errorf("request failed with status code %d", res.StatusCode)
		}
	}

	task := &struct {
		Task *models.Task `json:"task"`
	}{}

	err = json.NewDecoder(res.Body).Decode(&task)
	if err != nil {
		return nil, err
	}

	return task.Task, err
}

func (o *Orchestrator) PostResult(ctx context.Context, result *models.TaskResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", o.Url+"/internal/task", bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	res, err := o.Client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed status code %d", res.StatusCode)
	}

	return nil
}

func (o *Orchestrator) Ping() error {
	var err error
	var res *http.Response

	for i := 0; i <= o.Retries; i++ {
		time.Sleep(time.Duration(1<<i) * time.Second)
		res, err = http.Get(o.Url + "/internal/ping")
		if err == nil && res.StatusCode == http.StatusOK {
			res.Body.Close()
			return nil
		}
	}

	return err
}
