package request

type Login struct {
	Login    string `json:"login" validate:"required" example:"admin"`
	Password string `json:"password" validate:"min=4,max=50" example:"admin"`
}

type Refresh struct {
	RefreshToken string `json:"refresh_token" validate:"required" example:"refresh_token"`
}