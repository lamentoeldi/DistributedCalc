package models

type TaskResult struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Status string  `json:"status"`
	Final  bool    `json:"final"`
}

type AgentTask struct {
	Id            string  `json:"id"`
	LeftArg       float64 `json:"arg1"`
	RightArg      float64 `json:"arg2"`
	Op            string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
	Final         bool    `json:"final"`
}
