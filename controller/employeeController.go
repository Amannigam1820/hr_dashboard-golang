package controller

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Amannigam1820/hr-dashboard-golang/config"
	"github.com/Amannigam1820/hr-dashboard-golang/database"
	"github.com/Amannigam1820/hr-dashboard-golang/model"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateEmployee(c *fiber.Ctx) error {
	var employee model.Employee

	// Parse multipart form data (including file upload)
	err := c.BodyParser(&employee)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to parse request body")
	}

	// Handle file uploads
	form, err := c.MultipartForm()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Unable to parse multipart form")
	}

	// Process files (Resume, Experience Letter, Relieving Letter)
	if len(form.File["resume"]) > 0 {
		file := form.File["resume"][0] // *multipart.FileHeader
		f, err := file.Open()          // Open the file to get a multipart.File (the content of the file)
		if err != nil {
			log.Println("Error opening resume file:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to open resume file")
		}
		url, err := config.UploadToCloudinary(f) // Pass the file content to UploadToCloudinary
		if err != nil {
			log.Println("Error uploading resume:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload resume")
		}
		employee.Resume = url
		defer f.Close() // Don't forget to close the file after usage
	}

	if len(form.File["experience_letter"]) > 0 {
		file := form.File["experience_letter"][0] // *multipart.FileHeader
		f, err := file.Open()                     // Open the file to get a multipart.File (the content of the file)
		if err != nil {
			log.Println("Error opening experience letter file:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to open experience letter file")
		}
		url, err := config.UploadToCloudinary(f) // Pass the file content to UploadToCloudinary
		if err != nil {
			log.Println("Error uploading experience letter:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload experience letter")
		}
		employee.ExperienceLetter = url
		defer f.Close() // Don't forget to close the file after usage
	}

	if len(form.File["releiving_letter"]) > 0 {
		file := form.File["releiving_letter"][0] // *multipart.FileHeader
		f, err := file.Open()                    // Open the file to get a multipart.File (the content of the file)
		if err != nil {
			log.Println("Error opening relieving letter file:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to open relieving letter file")
		}
		url, err := config.UploadToCloudinary(f) // Pass the file content to UploadToCloudinary
		if err != nil {
			log.Println("Error uploading relieving letter:", err)
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload relieving letter")
		}
		employee.ReleivingLetter = url
		defer f.Close() // Don't forget to close the file after usage
	}

	// Set the created date
	employee.CreatedAt = time.Now()
	employee.UpdatedAt = time.Now()

	fmt.Println(&employee.CreatedAt)
	fmt.Println(&employee.BirthDate)

	// Insert employee into the database
	result := database.DBConn.Create(&employee)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": "false",
			"error":   "Failed to create HR record",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": "true",
		"message": "Hr Created SuccessFully",
		"Hr":      employee,
	})
}
func GetAllTechStackCategory(c *fiber.Ctx) error {
	var employees []model.Employee
	var techStacks []string

	// Fetch all tech stacks from the employees table
	if result := database.DBConn.Model(&model.Employee{}).Pluck("tech_stack", &employees); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve employee tech stacks",
		})
	}

	// Iterate over the employees and split the tech stacks by commas (for comma-separated values)
	for _, employee := range employees {
		// Assuming tech_stack is a comma-separated string like "MERN,Golang"
		techStacks = append(techStacks, strings.Split(employee.TechStack, ",")...)
	}

	// Create a map to track unique tech stacks
	uniqueTechStacks := make(map[string]bool)

	// Loop through the tech stacks and add them to the map (this will ensure uniqueness)
	var resultTechStacks []string
	for _, stack := range techStacks {
		if !uniqueTechStacks[stack] {
			uniqueTechStacks[stack] = true
			resultTechStacks = append(resultTechStacks, stack)
		}
	}

	// Return the response with the unique tech stacks
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"tech-stack": resultTechStacks,
		"success":    true,
	})
}
func GetEmployeeByTechStack(c *fiber.Ctx) error {
	// Retrieve the tech stack from the query parameters
	techStack := c.Query("techstack")
	if techStack == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "Tech stack query parameter is required",
		})
	}

	var employees []model.Employee
	// Find employees that match the tech stack in the tech_stack column
	if result := database.DBConn.Where("tech_stack LIKE ?", "%"+techStack+"%").Find(&employees); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve Employee records",
		})
	}

	// Return the filtered employee list
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    employees,
		"success": true,
	})
}

func GetMyProfile(c *fiber.Ctx) error {
	// Retrieve the logged-in user from the context
	user, ok := c.Locals("user").(model.Hr)
	fmt.Println("user :", user, "ok :", ok)

	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated or not found in context",
		})
	}

	// Query the employee table to find a matching record
	var employee model.Employee
	if err := database.DBConn.Where("email = ? AND name = ? ", user.Email, user.Name).First(&employee).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// If no matching record is found, return a 404 response
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "No matching employee found",
			})
		}
		// Handle other potential errors
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query employee table",
		})
	}

	// If a match is found, return all fields of the employee
	return c.JSON(employee)
}

func GetEmployeeStats(c *fiber.Ctx) error {
	var employees []model.Employee

	// Fetch all employees from the database
	if result := database.DBConn.Find(&employees); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve Employee records",
		})
	}

	// Initialize maps to store the distribution data
	genderDistribution := make(map[string]int)
	departmentDistribution := make(map[string]int)
	positionDistribution := make(map[string]int)
	locationDistribution := make(map[string]int)

	// Iterate through the employees to count distributions
	for _, employee := range employees {
		// Gender Distribution
		genderDistribution[employee.Gender]++

		// Department Distribution
		departmentDistribution[employee.Department]++

		// Postion Distribution
		positionDistribution[employee.Position]++

		// Location Distribution
		locationDistribution[employee.Address]++
	}

	// Return the aggregated data in the required format without age distribution
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"genderDistribution":     genderDistribution,
			"departmentDistribution": departmentDistribution,
			"positionDistribution":   positionDistribution,
			"locationDistribution":   locationDistribution,
		},
	})
}

func GetAllEmployee(c *fiber.Ctx) error {
	var employees []model.Employee

	if result := database.DBConn.Find(&employees); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to retrieve Employee records",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{

		"data":    employees,
		"success": true,
	})
}

func GetEmployeeById(c *fiber.Ctx) error {
	var employee model.Employee
	id := c.Params("id")

	if err := database.DBConn.First(&employee, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "HR record not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"hr":      employee,
	})
}

func DeleteEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee model.Employee

	if err := database.DBConn.First(&employee, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"error":   "Employee record not found",
		})
	}
	if result := database.DBConn.Delete(&employee); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to delete HR record",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Employee record deleted successfully",
	})
}

func UpdateEmployee(c *fiber.Ctx) error {
	id := c.Params("id")
	var employee model.Employee

	if err := c.BodyParser(&employee); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success":  false,
			"error":    "Failed to parse request body",
			"sytemErr": err,
		})
	}

	//fmt.Println(employee)

	var existingEmployee model.Employee
	if result := database.DBConn.First(&existingEmployee, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Employee Record not found",
		})
	}

	file, err := c.FormFile("resume")
	fmt.Println("filereror", file, err)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to get the file",
			"err":   err,
		})
	}
	fileContent, err := file.Open()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error opening file",
		})
	}
	defer fileContent.Close()

	// Upload the file to Cloudinary
	resumeURL, err := config.UploadToCloudinary(fileContent)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": fmt.Sprintf("Failed to upload to Cloudinary: %v", err),
		})
	}

	// Store the Cloudinary URL in the user struct
	existingEmployee.Resume = resumeURL
	if employee.Name != "" {
		existingEmployee.Name = employee.Name
	}
	if employee.ContactNumber != "" {
		existingEmployee.ContactNumber = employee.ContactNumber
	}
	if employee.Email != "" {
		existingEmployee.Email = employee.Email
		var existingEmail model.Employee
		if err := database.DBConn.Where("email:?", employee.Email).First(&existingEmail).Error; err == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "Email already exists",
			})
		}
	}
	if employee.TechStack != "" {
		existingEmployee.TechStack = employee.TechStack
	}
	if employee.Position != "" {
		existingEmployee.Position = employee.Position
	}
	if employee.YearsOfExperience != 0 {
		existingEmployee.YearsOfExperience = employee.YearsOfExperience
	}
	if employee.CasualLeave != 0 {
		existingEmployee.CasualLeave = employee.CasualLeave
	}
	if employee.EarnedLeave != 0 {
		existingEmployee.EarnedLeave = employee.EarnedLeave
	}
	if employee.Salary != 0 {
		existingEmployee.Salary = employee.Salary
	}
	if employee.Department != "" {
		existingEmployee.Department = employee.Department
	}
	if employee.Performance != "" {
		existingEmployee.Performance = employee.Performance
	}
	if employee.Address != "" {
		existingEmployee.Address = employee.Address
	}

	if result := database.DBConn.Save(&existingEmployee); result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Failed to update HR record",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "HR record updated successfully !",
		"data":    existingEmployee,
	})
}
