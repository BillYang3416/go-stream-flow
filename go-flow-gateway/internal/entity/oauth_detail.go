package entity

type OAuthDetail struct {
	OAuthID      string `json:"oauthId"`
	UserID       int    `json:"userId"`
	Provider     string `json:"provider"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
