package dto

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
	AccessToken  string `json:"access_token" binding:"required"`
}

type JWTTokenData struct {
	Token     string
	JTI       string
	ExpiresIn int
}

type LogoutRequest struct {
	AccessToken string `json:"access_token" binding:"required"`
}
