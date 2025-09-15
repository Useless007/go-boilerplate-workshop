package middleware

import (
	"strings"
	"boilerplate/utils"
	"github.com/gofiber/fiber/v2"
)

func JWTMiddleware() fiber.Handler {
	return func (c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Unauthorized",
			})
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token",
			})
		}
		c.Locals("user_id", claims.UserID)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

func RoleMiddleware(requireRole string) fiber.Handler {
	return func (c *fiber.Ctx) error {
		role := c.Locals("role").(string)
		if role != requireRole && role != "ADMIN" {
			return c.Status(403).JSON(fiber.Map{
				"success": false,
				"message": "Forbidden",
			})
		}
		return c.Next()
	}
}