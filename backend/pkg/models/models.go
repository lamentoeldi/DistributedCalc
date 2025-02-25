package models

type Task struct {
	Id            int     `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int64   `json:"operation_time"`
}

type TaskResult struct {
	Id     int     `json:"id"`
	Result float64 `json:"result"`
}

type Expression struct {
	Id     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type Calculation struct {
	Expression string `json:"expression"`
}
