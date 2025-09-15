package database

import (
	"boilerplate/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	db    *gorm.DB
	dbErr error
)

// DB Connect
func Connect() {
	db, dbErr = gorm.Open(sqlite.Open("users.db"), &gorm.Config{})
	if dbErr != nil {
		log.Fatal("Failed to connect to database: ", dbErr)
	}

	// Migrate the schema
	dbErr = db.AutoMigrate(&models.User{})
	if dbErr != nil {
		log.Fatal("Failed to migrate database:", dbErr)
	}

	log.Println("Database connection established")
}


func Insert(user *models.User) error {
	return db.Create(user).Error
}

func Get() []models.User {
	var users []models.User
	db.Find(&users)
	return users
}

func Update(id int, newName string) error {
	return db.Model(&models.User{}).Where("id = ?", id).Update("name", newName).Error
}

func Delete(id int) error {
	return db.Delete(&models.User{}, id).Error
}