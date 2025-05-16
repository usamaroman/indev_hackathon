package response

type Login struct {
	User         *User  `json:"user"`
	AccessToken  string `json:"access_token" example:"access_token"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
}

type Refresh struct {
	AccessToken  string `json:"access_token" example:"access_token"`
	RefreshToken string `json:"refresh_token" example:"refresh_token"`
}