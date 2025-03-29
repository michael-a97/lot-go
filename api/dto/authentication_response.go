package dto

type AuthenticationResponse struct {
	ID           uint   `json:"id"`
	Role         string `json:"role"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}
