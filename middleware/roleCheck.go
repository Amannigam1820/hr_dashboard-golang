package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func RoleCheck(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Retrieve token from the cookie or authorization header
		tokenString := c.Cookies("token")
		// fmt.Println(tokenString)
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token not provided",
			})
		}

		// clear=============================

		// Parse and verify the token
		secretKey := []byte("hrdashboard")
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			// Ensure signing method is correct
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return secretKey, nil
		})

		//	fmt.Println(token.Valid)
		//	fmt.Println(token, err)

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token 123",
			})
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		//fmt.Println(claims)

		// Check the user's role
		userRole := claims["role"].(string)
		for _, role := range roles {
			if userRole == role {
				return c.Next() // Role matches, proceed to the next handler
			}
		}

		// If role does not match, deny access
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Forbidden - You do not have the required permissions",
		})
	}
}
