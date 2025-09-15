package handlers

import (
	"boilerplate/database"
	"boilerplate/models"
	"boilerplate/utils"

	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")
	name := c.FormValue("name")
	role := c.FormValue("role", "USER") // default role is USER

	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password, // will be hashed in database.Insert
		Role:     role,
	}
	err := database.Insert(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

func Login(c *fiber.Ctx) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	user, err := database.GetUserByEmail(email)
	if err != nil || !utils.CheckPassword(user.Password, password) {
		return c.Status(401).JSON(fiber.Map{
			"success": false,
			"message": "Invalid email or password",
		})
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": "Could not generate token",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"token":   token,
		"user":    user,
	})
}

// UserList returns a list of users
func UserList(c *fiber.Ctx) error {
	users := database.Get()

	return c.JSON(fiber.Map{
		"success": true,
		"users":   users,
	})
}

// UserCreate registers a user
func UserCreate(c *fiber.Ctx) error {
	user := &models.User{
		// Note: when writing to external database,
		// we can simply use - Name: c.FormValue("user")
		Name: c.FormValue("user"),
	}
	err := database.Insert(user)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}

func UserUpdate(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid id",
		})
	}
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	if role != "ADMIN" && userID != id {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
		})
	}
	newName := c.FormValue("name")
	err = database.Update(id, newName)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
	})
}

func UserDelete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": "Invalid id",
		})
	}
	userID := c.Locals("user_id").(int)
	role := c.Locals("role").(string)
	if role != "ADMIN" && userID != id {
		return c.Status(403).JSON(fiber.Map{
			"success": false,
			"message": "Forbidden",
		})
	}
	err = database.Delete(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
	})
}

// NotFound returns custom 404 page
func NotFound(c *fiber.Ctx) error {
	return c.Status(404).SendFile("./static/private/404.html")
}

func Me(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(int)
	user, err := database.GetUserByID(userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"success": false,
			"message": "User not found",
		})
	}
	return c.JSON(fiber.Map{
		"success": true,
		"user":    user,
	})
}
