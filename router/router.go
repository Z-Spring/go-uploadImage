package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"upload2/global"
	"upload2/upload"
)

func NewRouter() {
	router := gin.Default()
	router.StaticFS("/static", http.Dir(global.AppSetting.UploadStoreUrl))
	router.POST("/upload", UploadImage)
	router.Run(":8080")
}

func UploadImage(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("error?", err)
	}
	fileInfo, err := upload.ImageUpload(fileHeader, file)
	if err != nil {
		log.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"fileUrl": fileInfo.FileUrl})
}
