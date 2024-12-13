package middleware

import (
	"fmt"

	"github.com/Amannigam1820/hr-dashboard-golang/database"
	"github.com/Amannigam1820/hr-dashboard-golang/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("hrdashboard") // Replace with your actual secret key

func IsAuthenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get token from cookies
		tokenString := c.Cookies("token")
		//fmt.Println("Token from cookie:", tokenString)

		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication token is missing",
			})
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		//fmt.Println("token : ", token)
		if err != nil || !token.Valid {
			//fmt.Println("Token parsing failed:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// everything right in this place

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		//fmt.Println("claims :", claims, "ok:", ok)

		if !ok || claims["id"] == nil {
			//	fmt.Println("Invalid token claims:", claims)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// Extract user_id from token claims
		userID := uint(claims["id"].(float64))
		//fmt.Println("Extracted user_id from token claims:", claims["id"])

		// Fetch user details from the database
		var user model.Hr
		err = database.DBConn.First(&user, "id = ?", userID).Error
		if err != nil {
			//	fmt.Println("User not found in DB:", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not found",
			})
		}

		// Store user in context
		//fmt.Println("Authenticated user:", user)
		c.Locals("user", user)

		return c.Next()
	}
}
