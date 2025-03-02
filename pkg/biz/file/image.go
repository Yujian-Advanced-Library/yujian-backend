package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"yujian-backend/pkg/file"
	"yujian-backend/pkg/model"
)

func FetchFile() gin.HandlerFunc {
	return func(c *gin.Context) {
		imageId := c.Param("imageId")
		client := file.GetMinioClient()
		fileName := imageId + ".jpg"
		if fetchFile, err := client.DownloadFile(c, "images", imageId, fileName); err != nil {
			c.JSON(http.StatusInternalServerError, model.BaseResp{Code: http.StatusInternalServerError, ErrMsg: "failed to fetch file", Error: err})
			return
		} else {
			c.Header("Content-Type", "image/jpeg")
			c.File(fetchFile)
		}
	}
}
