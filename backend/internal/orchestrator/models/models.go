package models

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
