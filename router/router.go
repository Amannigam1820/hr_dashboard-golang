package router

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {

	SetupHRRoutes(app)
	SetupEmployeeRoutes(app)
	//SetupHRRoutes(app)

}
