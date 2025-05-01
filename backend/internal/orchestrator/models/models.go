package models

type Task struct {
	ID       string  `bson:"_id"`
	ExpID    string  `bson:"exp_id"`
	Op       string  `bson:"op"`
	LeftID   *string `bson:"left_id"`
	RightID  *string `bson:"right_id"`
	LeftArg  float64 `bson:"left_arg"`
	RightArg float64 `bson:"right_arg"`
	Result   float64 `bson:"result"`
	Status   string  `bson:"status"`
	Final    bool    `bson:"final"`
}

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

type Expression struct {
	Id     string  `json:"id"`
	Result float64 `json:"result"`
	Status string  `json:"status"`
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}
