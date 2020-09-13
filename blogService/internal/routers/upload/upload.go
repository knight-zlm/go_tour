package upload

import (
	"github.com/gin-gonic/gin"
	"github.com/knight-zlm/blog-service/internal/service"
	"github.com/knight-zlm/blog-service/pkg/convert"
	"github.com/knight-zlm/blog-service/pkg/errcode"
	"github.com/knight-zlm/blog-service/pkg/upload"

	"github.com/knight-zlm/blog-service/pkg/app"
)

func UploadFile(c *gin.Context) {
	response := app.NewResponse(c)
	file, fileHeader, err := c.Request.FormFile("file")
	fileType := convert.StrTo(c.PostForm("type")).MustInt()
	if err != nil {
		response.ToErrorResponse(errcode.InvalidParams.WithDetails(err.Error()))
		return
	}
	if fileHeader == nil || fileType < 0 {

	}

	svc := service.New(c)
	fileInfo, err := svc.UploadFile(upload.FileType(fileType), file, fileHeader)
	if err != nil {
		response.ToErrorResponse(errcode.ErrorUploadFileFail.WithDetails(err.Error()))
		return
	}

	response.ToResponse(gin.H{"file_access_url": fileInfo.AccessUrl})
}
