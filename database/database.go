package database

import (
	"boilerplate/models"
	"fmt"
	"sync"
)

var (
	db []*models.User
	mu sync.Mutex
	nextID = 0
)

// Connect with database
func Connect() {
	db = make([]*models.User, 0)
	fmt.Println("Connected with Database")
}

func Insert(user *models.User) { //create
	mu.Lock()
	defer mu.Unlock()
	user.ID = nextID
	nextID++
	db = append(db, user)
}

func Get() []*models.User { //read
	return db
}

func Update(id int, newName string) error {
	mu.Lock()
	defer mu.Unlock()

	for _, user := range db {
		if user.ID == id {
			user.Name = newName
			return nil
		}
	}

	return fmt.Errorf("user with id %d not found", id)
}

func Delete(id int) error {
	mu.Lock()
	defer mu.Unlock()

	for i, user := range db {
		if user.ID == id {
			db = append(db[:i], db[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("user with id %d not found", id)
}
