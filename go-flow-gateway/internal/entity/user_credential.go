package entity

type UserCredential struct {
	CredentialID string `json:"credentialId"`
	UserID       int    `json:"userId"`
	Username     string `json:"username"`
	PasswordHash string `json:"password"`
}
