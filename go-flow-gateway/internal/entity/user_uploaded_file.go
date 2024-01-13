package entity

import "time"

// File represents the file-related information that will be stored and retrieved.
type UserUploadedFile struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	Size           int64      `json:"size"`
	Content        []byte     `json:"-"`
	Base64Content  string     `json:"content,omitempty"`
	UserID         string     `json:"userId"`
	CreatedAt      time.Time  `json:"createdAt"`
	EmailSent      bool       `json:"emailSent"`      // Indicates if the email was sent successfully
	EmailSentAt    *time.Time `json:"emailSentAt"`    // The timestamp when the email was sent
	EmailRecipient string     `json:"emailRecipient"` // The email address of the recipient
	ErrorMessage   string     `json:"errorMessage"`   //  // Error message if the email was not sent successfully
}
