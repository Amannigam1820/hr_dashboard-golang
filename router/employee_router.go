package router

import (
	"github.com/Amannigam1820/hr-dashboard-golang/controller"
	"github.com/Amannigam1820/hr-dashboard-golang/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupEmployeeRoutes(app *fiber.App) {
	employeeGroup := app.Group("/api/employee")

	employeeGroup.Post("/", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.CreateEmployee)
	employeeGroup.Get("/all", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.GetAllEmployee)
	employeeGroup.Get("/stats", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.GetEmployeeStats)
	employeeGroup.Get("/:id", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.GetEmployeeById)

	employeeGroup.Delete("/:id", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.DeleteEmployee)
	employeeGroup.Put("/:id", middleware.RoleCheck("Hr-Admin", "Super-Admin"), controller.UpdateEmployee)

}
