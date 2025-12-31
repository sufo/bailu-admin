/**
 * Create by sufo
 * @Email ouamour@gmail.com
 *
 * @Desc 上传处理
 */

package system

import (
	"bailu/app/api/admin"
	"bailu/app/domain/resp"
	"bailu/global/consts"
	respErr "bailu/pkg/exception"
	"bailu/utils/upload"
	"github.com/gin-gonic/gin"
)

type UploadApi struct {
	Oss upload.OSS
}

func NewUploadApi(Oss upload.OSS) *UploadApi {
	return &UploadApi{Oss}
}

// upload
// @title 上传
// @Summary 上传接口
// @Description 上传接口
// @Tags Upload
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "file"
// @Success 200 {object} resp.Response[any]
// @Router /api/post [post]
// @Security Bearer
func (u *UploadApi) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil || file == nil {
		respErr.BadRequestErrorWithError(err)
	}
	fPath, _, err := u.Oss.UploadFileToDir(file, consts.IMG_DIR)
	if err != nil {
		respErr.InternalServerErrorWithError(err)
	}
	r := c.Request
	url := admin.FileUrl(r, fPath)
	resp.OKWithData(c, map[string]string{
		"url": url,
	})
}
