package models

// User model
type User struct {
	ID    int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name  string `json:"name"`
	Email string `json:"email" gorm:"unique"`
	Password string `json:"-"` // omit password in JSON responses
	Role  string `json:"role" gorm:"default:USER"` // default role is 'USER'
}
