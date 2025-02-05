package files

import (
	"context"
	"io/ioutil"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	request "github.com/wafi04/backend/pkg/types/req"
	response "github.com/wafi04/backend/pkg/types/res"
)

type Cloudinary struct {
	cloudinary *cloudinary.Cloudinary
}

func NewCloudinaryService(cld *cloudinary.Cloudinary) *Cloudinary {
	return &Cloudinary{cloudinary: cld}
}

func (s *Cloudinary) UploadFile(
	ctx context.Context,
	req *request.FileUploadRequest,
) (*response.FileUploadResponse, error) {
	tempFile, err := ioutil.TempFile("", "upload-*")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tempFile.Name())

	// Write the multipart file to temp file
	fileBytes, err := ioutil.ReadAll(req.FileData)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(tempFile.Name(), fileBytes, 0644); err != nil {
		return nil, err
	}

	tempFile.Close()

	uploadResult, err := s.cloudinary.Upload.Upload(ctx, tempFile.Name(), uploader.UploadParams{
		Folder:   req.Folder,
		PublicID: req.PublicID,
	})
	if err != nil {
		return nil, err
	}

	return &response.FileUploadResponse{
		URL:      uploadResult.SecureURL,
		PublicID: uploadResult.PublicID,
	}, nil
}
