package config

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

var CloudinaryClient *cloudinary.Cloudinary

func InitCloudinary() {
	cld, err := cloudinary.NewFromParams("dcohtcgb7", "997959291724572", "RYS6MeK3gp5qCG3OfSj99j_wwuw")
	if err != nil {
		log.Fatal("Error initializing Cloudinary: ", err)
	}
	CloudinaryClient = cld
}

func UploadToCloudinary(file multipart.File) (string, error) {
	if CloudinaryClient == nil {
		return "", fmt.Errorf("cloudinary client is not initialized")
	}

	// Generate a unique public ID by adding a UUID or timestamp
	uniqueID := uuid.New().String()
	// Upload to Cloudinary using the latest API
	uploadResult, err := CloudinaryClient.Upload.Upload(context.Background(), file, uploader.UploadParams{
		Folder:       "employee_resume_pdf",
		PublicID:     uniqueID,
		ResourceType: "raw",
		//AccessControlAllowOrigin:"*",

	})
	if err != nil {
		return "", fmt.Errorf("error uploading to Cloudinary: %v", err)
	}
	return uploadResult.SecureURL, nil
}

// package config

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"mime/multipart"
// 	"net/http"
// 	"strings"

// 	"github.com/cloudinary/cloudinary-go/v2"
// 	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
// )

// var CloudinaryClient *cloudinary.Cloudinary

// func InitCloudinary() {
// 	cld, err := cloudinary.NewFromParams("dcohtcgb7", "997959291724572", "RYS6MeK3gp5qCG3OfSj99j_wwuw")
// 	if err != nil {
// 		log.Fatal("Error initializing Cloudinary: ", err)
// 	}
// 	CloudinaryClient = cld
// }

// func UploadToCloudinary(file multipart.File) (string, error) {
// 	if CloudinaryClient == nil {
// 		return "", fmt.Errorf("Cloudinary client is not initialized")
// 	}

// 	// Check the MIME type of the file
// 	buffer := make([]byte, 512)
// 	_, err := file.Read(buffer)
// 	if err != nil {
// 		return "", fmt.Errorf("error reading file: %v", err)
// 	}
// 	fileType := http.DetectContentType(buffer)

// 	// Check if the file is a PDF
// 	if !strings.HasPrefix(fileType, "application/pdf") {
// 		return "", fmt.Errorf("the file is not a PDF")
// 	}

// 	// Reset file pointer before upload
// 	_, err = file.Seek(0, 0)
// 	if err != nil {
// 		return "", fmt.Errorf("error resetting file pointer: %v", err)
// 	}

// 	// Upload to Cloudinary with a folder name (optional)
// 	uploadResult, err := CloudinaryClient.Upload.Upload(context.Background(), file, uploader.UploadParams{
// 		Folder: "employee_resume_pdf", // Optional folder name in Cloudinary
// 		ResourceType:         "auto",                 // Specify raw file type (for PDF)
//     AccessControlAllowOrigin: "*",
// 	})
// 	if err != nil {
// 		return "", fmt.Errorf("error uploading to Cloudinary: %v", err)
// 	}

// 	fmt.Println(":upload", uploadResult)

// 	// Return the URL of the uploaded PDF
// 	return uploadResult.SecureURL, nil
// }
