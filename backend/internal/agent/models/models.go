package models

type TaskResult struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Status string  `json:"status"`
	Final  bool    `json:"final"`
}

type AgentTask struct {
	Id            string  `json:"id"`
	LeftArg       float64 `json:"left_arg"`
	RightArg      float64 `json:"right_arg"`
	Op            string  `json:"op"`
	OperationTime int64   `json:"op_time"`
	Final         bool    `json:"final"`
}
