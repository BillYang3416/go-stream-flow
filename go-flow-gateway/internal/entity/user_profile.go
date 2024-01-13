package entity

type UserProfile struct {
	UserID       string `json:"userId"`
	DisplayName  string `json:"displayName"`
	PictureURL   string `json:"pictureUrl"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
