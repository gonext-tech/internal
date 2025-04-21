package services

import (
	"errors"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"gorm.io/gorm"
)

const (
	maxImageWidth  = 800
	maxImageHeight = 600
)

var allowedExtensions = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".webp": true,
}

type UploadServices struct {
	DB *gorm.DB
}

func NewUploadServices(db *gorm.DB) *UploadServices {
	return &UploadServices{
		DB: db,
	}
}

func (us *UploadServices) Upload(file *multipart.FileHeader, project string, folder string) (string, error) {
	//Check the extension of the image
	if !us.isValidImageExtension(file.Filename) {
		return "", errors.New("image extension is not available")
	}
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}

	defer src.Close()

	// Decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return "", err
	}

	// Resize image if necessary
	img = us.resizeImage(img)

	folderURL := os.Getenv("BUCKET_FOLDER")
	homefolderPath, err := homedir.Expand(folderURL)
	if err != nil {
		return "", err
	}

	// Generate unique file name
	fileName := us.generateFileName(file.Filename)

	// Create directory if not exists
	dir := filepath.Join(homefolderPath, strings.ToLower(project), folder)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", err
		}
	}

	// Create destination file
	dstPath := filepath.Join(dir, fileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	// Save resized image to destination
	if err := webp.Encode(dst, img, &webp.Options{Quality: 80}); err != nil {
		return "", err
	}
	// Generate URL for the uploaded image

	bucketURL := os.Getenv("BUCKET_URL")
	imageURL := fmt.Sprintf("%s/%s/%s/%s", bucketURL, strings.ToLower(project), folder, fileName)

	return imageURL, nil

}

func (us *UploadServices) Delete(imageName string) error {
	if !us.isValidImageExtension(imageName) {
		return errors.New("image extension is not available")
	}

	folderURL := os.Getenv("BUCKET_FOLDER")
	homefolderPath, err := homedir.Expand(folderURL)
	if err != nil {
		return err
	}

	bucketURL := os.Getenv("BUCKET_URL")
	image := strings.TrimPrefix(imageName, bucketURL)
	// Construct path to the image file
	imagePath := filepath.Join(homefolderPath, image)

	// Check if the image file exists
	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		return errors.New("image not found")
	}

	// Delete the image file
	if err := os.Remove(imagePath); err != nil {
		return errors.New("cannot delete the image")
	}
	return nil
}

func (us *UploadServices) Serve(imagePath string) (io.ReadCloser, error) {
	if !us.isValidImageExtension(imagePath) {
		return nil, errors.New("image extension is not available")
	}
	folderURL := os.Getenv("BUCKET_FOLDER")
	homefolderPath, err := homedir.Expand(folderURL)
	if err != nil {
		return nil, err
	}
	filePath := filepath.Join(homefolderPath, imagePath)

	// Check if the file exists
	_, err = os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			// File not found
			return nil, errors.New("file not found")
		}
		// Other error
		return nil, errors.New("internal server error")
	}
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New("file not found")
	}

	// Return a reader that starts from the beginning of the file
	return file, nil
}

func (us *UploadServices) resizeImage(img image.Image) image.Image {
	// Resize image if necessary
	if img.Bounds().Dx() > maxImageWidth || img.Bounds().Dy() > maxImageHeight {
		img = imaging.Resize(img, maxImageWidth, maxImageHeight, imaging.Lanczos)
	}
	return img
}

func (us *UploadServices) generateFileName(originalName string) string {
	// Generate UUID
	uuid := uuid.New()

	// Generate timestamp
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Construct file name
	return fmt.Sprintf("%s_%d%s", uuid, timestamp, filepath.Ext(originalName))
}

func (us *UploadServices) isValidImageExtension(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return allowedExtensions[ext]
}
