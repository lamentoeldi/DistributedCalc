package models

type UserCredentials struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type JWTTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
