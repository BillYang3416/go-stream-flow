package entity

type UserProfile struct {
	UserID      int    `json:"userId"`
	DisplayName string `json:"displayName"`
	PictureURL  string `json:"pictureUrl"`
}
