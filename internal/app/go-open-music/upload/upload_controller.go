package upload

import (
	"net/http"
	"os"
	"strings"

	"github.com/TesyarRAz/go-open-music/internal/pkg/config"
	"github.com/TesyarRAz/go-open-music/internal/pkg/util"
	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func (UploadController) Store(c *gin.Context) {
	file, err := c.FormFile("data")

	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if file.Size/1000 > 500 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "file lebih dari 500kb",
		})
		return
	}

	ext := file.Filename[strings.LastIndex(file.Filename, "."):]

	os.MkdirAll("static/images/", os.ModePerm)
	filePath := "static/images/" + util.RandomText(10) + ext
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	if ext != ".jpg" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "File bukan gambar",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Berhasil upload gambar",
		"data": gin.H{
			"pictureUrl": strings.TrimSuffix(config.AppConfig.APP_URL, "/") + "/" + filePath,
		},
	})
}
