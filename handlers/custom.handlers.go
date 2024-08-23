package handlers

import (
	"log"

	"github.com/labstack/echo/v4"
)

func UploadImage(c echo.Context, us UploadService, project, folder string) []string {
	var imageURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		log.Println("err-1", err)
		return imageURLs
	}

	files := form.File["file"]
	if len(files) == 0 {
		return imageURLs
	}
	oldImage := c.FormValue("image")
	if oldImage != "" {
		err = us.Delete(oldImage)
		if err != nil {
			log.Println("errr-delete", err)
		}
	}

	for _, file := range files {
		imageURL, err := us.Upload(file, project, folder)
		if err != nil {
			setFlashmessages(c, "error", err.Error())
			return imageURLs
		}
		imageURLs = append(imageURLs, imageURL)
	}

	return imageURLs
}
