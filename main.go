package main

import (
	"github.com/Amannigam1820/hr-dashboard-golang/config"
	"github.com/Amannigam1820/hr-dashboard-golang/database"
	"github.com/Amannigam1820/hr-dashboard-golang/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	database.ConnectDB()
}

func main() {
	config.InitCloudinary()

	sqlDb, err := database.DBConn.DB()
	if err != nil {
		panic("Error in sql connection")
	}
	defer sqlDb.Close()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173", // Replace with your frontend URL
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: true, // Allow cookies/credentials
	}))
	// app.Use(cors.New(cors.Config{
	// 	AllowCredentials: true,  // Allow sending cookies with requests
	// 	AllowOrigins:     "http://localhost:5173",  // Allow the frontend origin
	// 	AllowMethods:     "GET,POST,PUT,DELETE",  // Allow HTTP methods
	// 	AllowHeaders:     "Content-Type, Authorization",  // Allow specific headers
	// }))

	router.SetupRoutes(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": " qwerty Welcome to my first api in fiber"})

	})

	app.Listen(":8080")
}
