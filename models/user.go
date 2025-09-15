package models

// User model
type User struct {
	ID  int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name"`
}
