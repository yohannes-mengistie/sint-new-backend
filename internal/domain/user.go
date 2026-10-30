package domain

type User struct {
	FirebaseUID string `json:"firebase_uid"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}
