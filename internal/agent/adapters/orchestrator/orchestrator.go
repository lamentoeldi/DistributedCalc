package orchestrator

import (
	"DistributedCalc/pkg/models"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
)

type Orchestrator struct {
	Client *http.Client
	Url    string // protocol://host:port
}

func NewOrchestrator(client *http.Client, url string) *Orchestrator {
	return &Orchestrator{
		Client: client,
		Url:    url,
	}
}

func (o *Orchestrator) GetTask(_ context.Context) (*models.Task, error) {
	req, err := http.NewRequest("GET", o.Url+"/internal/task", nil)
	if err != nil {
		return nil, err
	}

	res, err := o.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var task *models.Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		return nil, err
	}

	return task, err
}

func (o *Orchestrator) PostResult(_ context.Context, result *models.TaskResult) error {
	data, err := json.Marshal(result)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", o.Url+"/internal/task", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer req.Body.Close()

	_, err = o.Client.Do(req)
	if err != nil {
		return err
	}

	return nil
}
