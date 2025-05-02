package models

type Task struct {
	ID       string  `bson:"_id"`
	ExpID    string  `bson:"exp_id"`
	Op       string  `bson:"op"`
	LeftID   *string `bson:"left_id,omitempty"`
	RightID  *string `bson:"right_id,omitempty"`
	LeftArg  float64 `bson:"left_arg"`
	RightArg float64 `bson:"right_arg"`
	Result   float64
	Status   string `bson:"status"`
	Final    bool   `bson:"final"`
}

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

type Expression struct {
	Id     string  `json:"id" bson:"_id"`
	UserID string  `json:"user_id" bson:"user_id"`
	Result float64 `json:"result" bson:"result"`
	Status string  `json:"status" bson:"status"`
}

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type UserCredentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type User struct {
	Id             string `json:"id" bson:"_id"`
	Username       string `json:"username" bson:"username"`
	HashedPassword []byte `json:"hashed_password" bson:"hashed_password"`
}
