package router

import (
	"github.com/gin-gonic/gin"
	"github.com/winstonajr/dashdrop/internal/handler"
)

func SetupRouter(fileHandler *handler.FileHandler) *gin.Engine {
	r := gin.Default()
	r.StaticFile("/", "./public/index.html")

	api := r.Group("/api")
	{
		api.GET("/files", fileHandler.ListFiles)
		api.POST("/upload", fileHandler.UploadFiles)
		api.GET("/download/:id", fileHandler.DownloadFile)
		api.DELETE("/files/:id", fileHandler.DeleteFile)
	}

	return r
}
