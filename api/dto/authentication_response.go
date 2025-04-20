package dto

type AuthenticationResponse struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	User         *UserDto `json:"user"`
}
