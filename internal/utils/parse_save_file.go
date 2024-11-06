package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func ParseAndSaveFiles(r *http.Request, productID string) ([]string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return nil, err
	}

	productDir := "/app/static/products/"
	var imageURLs []string
	files := r.MultipartForm.File["images"]

	for i, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer file.Close()

		filePath := fmt.Sprintf("%s%s-%d.webp", productDir, productID, i+1)
		dst, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file: %v\n", err)
			return nil, err
		}
		defer dst.Close()

		fmt.Printf("Saving file to: %s\n", filePath)

		if _, err = io.Copy(dst, file); err != nil {
			fmt.Printf("Error copying file: %v\n", err)
			return nil, err
		}

		BaseUrlToProducts, ok := os.LookupEnv("BASE_URL_TO_PRODUCTS")
		if !ok {
			return nil, fmt.Errorf("BASE_URL_TO_PRODUCTS environment variable not found")
		}

		imageURL := fmt.Sprintf("%s/%s-%d.webp", BaseUrlToProducts, productID, i+1)
		imageURLs = append(imageURLs, imageURL)
	}

	return imageURLs, nil
}

func ParseFormValueAsInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}
