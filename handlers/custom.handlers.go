package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func UploadImage(c echo.Context, us UploadService, project, folder string) []string {
	var imageURLs []string
	form, err := c.MultipartForm()
	if err != nil {
		log.Error(err)
		return imageURLs
	}

	files := form.File["file"]
	if len(files) == 0 {
		log.Error(errors.New("theres not file"))
		return imageURLs
	}
	oldImage := c.FormValue("image")
	if oldImage != "" {
		err = us.Delete(oldImage)
		if err != nil {
			return imageURLs
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
