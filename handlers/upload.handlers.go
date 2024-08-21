package handlers

import (
	"io"
	"mime/multipart"
)

type UploadService interface {
	Upload(file *multipart.FileHeader, project string, folder string) (string, error)
	Delete(imageName string) error
	Serve(imageName string) (io.ReadCloser, error)
}

type UploadHander struct {
	ProjectServices ProjectService
	UploadServices  UploadService
}

func NewUploadHandler(us UploadService, ps ProjectService) *UploadHander {
	return &UploadHander{
		ProjectServices: ps,
		UploadServices:  us,
	}
}
