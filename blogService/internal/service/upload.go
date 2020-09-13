package service

import (
	"errors"
	"mime/multipart"
	"os"
	"path"

	"github.com/knight-zlm/blog-service/global"

	"github.com/knight-zlm/blog-service/pkg/upload"
)

type FileInfo struct {
	Name      string
	AccessUrl string
}

func (s *Service) UploadFile(fileType upload.FileType, file multipart.File,
	fileHeader *multipart.FileHeader) (*FileInfo, error) {
	fileName := upload.GetFileName(fileHeader.Filename)
	uploadSavePath := upload.GetSavePath()
	dst := path.Join(uploadSavePath, fileName)
	if !upload.CheckContainExt(fileType, fileName) {
		return nil, errors.New("file suffix is not supported.")
	}

	if upload.CheckSavePath(dst) {
		if err := upload.CreateSavePath(dst, os.ModePerm); err != nil {
			return nil, errors.New("failed to create save directory.")
		}
	}

	if upload.CheckMaxSize(fileType, file) {
		return nil, errors.New("exceeded maximum file limit.")
	}

	if upload.CheckPermission(dst) {
		return nil, errors.New("insufficient file permissions.")
	}

	if err := upload.SaveFile(fileHeader, dst); err != nil {
		return nil, err
	}

	accessUrl := global.AppSetting.UploadServerUrl + "/" + fileName
	return &FileInfo{Name: fileName, AccessUrl: accessUrl}, nil
}
