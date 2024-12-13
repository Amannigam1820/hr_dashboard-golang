package controller

import (
	"fmt"
	"log"
	"time"

	"github.com/Amannigam1820/hr-dashboard-golang/config"
	"github.com/Amannigam1820/hr-dashboard-golang/database"
	"github.com/Amannigam1820/hr-dashboard-golang/model"
	"github.com/gofiber/fiber/v2"
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
	var existingEmployee model.Employee
	if result := database.DBConn.First(&existingEmployee, id); result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Employee Record not found",
		})
	}
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

	//fmt.Println(&existingEmployee.YearsOfExperience)

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

// func getUserResume(c *fiber.Ctx) error{
// 	var employee model.Employee
// 	id := c.Params("id")

// 	if err := database.DBConn.First(&employee, id).Error; err != nil {
// 		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
// 			"success": false,
// 			"error":   "HR record not found",
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"success": true,
// 		"hr":      employee,
// 	})
// }

// func UpdateEmployee(c *fiber.Ctx) error {
//     id := c.Params("id")
//     var employee model.Employee

//     // Parse the request body
//     if err := c.BodyParser(&employee); err != nil {
//         return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//             "success": false,
//             "error":   "Failed to parse request body",
//             "systemErr": err,
//         })
//     }

//     // Fetch the existing employee record
//     var existingEmployee model.Employee
//     if result := database.DBConn.First(&existingEmployee, id); result.Error != nil {
//         return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
//             "success": false,
//             "message": "Employee record not found",
//         })
//     }

//     // Update employee fields
//     if employee.Name != "" {
//         existingEmployee.Name = employee.Name
//     }
//     if employee.ContactNumber != "" {
//         existingEmployee.ContactNumber = employee.ContactNumber
//     }
//     if employee.Email != "" {
//         existingEmployee.Email = employee.Email
//         var existingEmail model.Employee
//         if err := database.DBConn.Where("email = ?", employee.Email).First(&existingEmail).Error; err == nil {
//             return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//                 "success": false,
//                 "error":   "Email already exists",
//             })
//         }
//     }
//     if employee.TechStack != "" {
//         existingEmployee.TechStack = employee.TechStack
//     }
//     if employee.Position != "" {
//         existingEmployee.Position = employee.Position
//     }
//     if employee.YearsOfExperience != 0 {
//         existingEmployee.YearsOfExperience = employee.YearsOfExperience
//     }
//     if employee.CasualLeave != 0 {
//         existingEmployee.CasualLeave = employee.CasualLeave
//     }
//     if employee.EarnedLeave != 0 {
//         existingEmployee.EarnedLeave = employee.EarnedLeave
//     }
//     if employee.Salary != 0 {
//         existingEmployee.Salary = employee.Salary
//     }
//     if employee.Department != "" {
//         existingEmployee.Department = employee.Department
//     }
//     if employee.Performance != "" {
//         existingEmployee.Performance = employee.Performance
//     }
//     if employee.Address != "" {
//         existingEmployee.Address = employee.Address
//     }

// 	form, err := c.MultipartForm()
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusBadRequest, "Unable to parse multipart form")
// 	}

// 	// Process files (Resume, Experience Letter, Relieving Letter)
// 	if len(form.File["resume"]) > 0 {
// 		file := form.File["resume"][0] // *multipart.FileHeader
// 		f, err := file.Open()          // Open the file to get a multipart.File (the content of the file)
// 		if err != nil {

// 			log.Println("Error opening resume file:", err)
// 			return fiber.NewError(fiber.StatusInternalServerError, "Failed to open resume file")
// 		}
// 		url, err := config.UploadToCloudinary(f) // Pass the file content to UploadToCloudinary
// 		if err != nil {
// 			log.Println("Error uploading resume:", err)
// 			return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload resume")
// 		}
// 		fmt.Println("url resume:", url)
// 		existingEmployee.Resume = url
// 		defer f.Close() // Don't forget to close the file after usage
// 	}

//     // Handle file uploads
//     // form, err := c.MultipartForm()
// 	// if err != nil {
// 	// 	return fiber.NewError(fiber.StatusBadRequest, "Unable to parse multipart form")
// 	// }
//     // if err == nil {
//     //     // Process Resume
//     //     if len(form.File["resume"]) > 0 {
//     //         file := form.File["resume"][0]
//     //         f, err := file.Open()
//     //         if err != nil {
//     //             log.Println("Error opening resume file:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to open resume file")
//     //         }
//     //         url, err := config.UploadToCloudinary(f)
//     //         if err != nil {
//     //             log.Println("Error uploading resume:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload resume")
//     //         }
//     //         fmt.Println("resume url", url)
//     //         existingEmployee.Resume = url
//     //         defer f.Close()
//     //     }

//     //     // Process Experience Letter
//     //     if len(form.File["experience_letter"]) > 0 {
//     //         file := form.File["experience_letter"][0]
//     //         f, err := file.Open()
//     //         if err != nil {
//     //             log.Println("Error opening experience letter file:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to open experience letter file")
//     //         }
//     //         url, err := config.UploadToCloudinary(f)
//     //         if err != nil {
//     //             log.Println("Error uploading experience letter:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload experience letter")
//     //         }
//     //         existingEmployee.ExperienceLetter = url
//     //         defer f.Close()
//     //     }

//     //     // Process Relieving Letter
//     //     if len(form.File["releiving_letter"]) > 0 {
//     //         file := form.File["releiving_letter"][0]
//     //         f, err := file.Open()
//     //         if err != nil {
//     //             log.Println("Error opening relieving letter file:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to open relieving letter file")
//     //         }
//     //         url, err := config.UploadToCloudinary(f)
//     //         if err != nil {
//     //             log.Println("Error uploading relieving letter:", err)
//     //             return fiber.NewError(fiber.StatusInternalServerError, "Failed to upload relieving letter")
//     //         }
//     //         existingEmployee.ReleivingLetter = url
//     //         defer f.Close()
//     //     }
//     // }

//     // Update the updated_at timestamp
//     existingEmployee.UpdatedAt = time.Now()

//     // Save the updated employee record
//     if result := database.DBConn.Save(&existingEmployee); result.Error != nil {
//         return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
//             "success": false,
//             "error":   "Failed to update employee record",
//         })
//     }

//     return c.Status(fiber.StatusOK).JSON(fiber.Map{
//         "success": true,
//         "message": "Employee record updated successfully!",
//         "data":    existingEmployee,
//     })
// }
