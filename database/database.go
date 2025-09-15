package database

import (
	"boilerplate/models"
	"boilerplate/utils"
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
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed
	return db.Create(user).Error
}

func Get() []models.User {
	var users []models.User
	db.Find(&users)
	return users
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func GetUserByID(id int) (*models.User, error) {
	var user models.User
	err := db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func Update(id int, newName string) error {
	return db.Model(&models.User{}).Where("id = ?", id).Update("name", newName).Error
}

func Delete(id int) error {
	return db.Delete(&models.User{}, id).Error
}
