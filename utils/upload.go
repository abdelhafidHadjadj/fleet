package utils

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

func SaveFile(c *fiber.Ctx, fieldName, saveDir string) (string, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return "", err // No file uploaded
		}
		log.Printf("Error retrieving file: %v", err) // Log error
		return "", err
	}

	// Create the directory if it doesn't exist
	if _, err := os.Stat(saveDir); os.IsNotExist(err) {
		err = os.MkdirAll(saveDir, os.ModePerm)
		if err != nil {
			log.Printf("Error creating directory: %v", err) // Log error
			return "", err
		}
	}

	// Save the file with a unique name
	filePath := filepath.Join(saveDir, file.Filename)
	err = c.SaveFile(file, filePath)
	if err != nil {
		log.Printf("Error saving file: %v", err) // Log error
		return "", err
	}

	log.Printf("File saved at: %s", filePath) // Log success
	return filePath, nil
}
