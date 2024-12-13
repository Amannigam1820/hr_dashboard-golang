package database

import (
	"log"

	"github.com/Amannigam1820/hr-dashboard-golang/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DBConn *gorm.DB

func ConnectDB() {
	dsn := "root:Aman12345@tcp(127.0.0.1:3306)/hr_dashboard?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})

	if err != nil {
		panic("Database connection failed.")
	}
	log.Println("Database connection successfully")
	err = db.AutoMigrate(
		&model.Hr{},
		&model.Employee{},
	)
	if err != nil {
		log.Println("Some error ")
	}
	DBConn = db
}
