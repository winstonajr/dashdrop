package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/winstonajr/dashdrop/internal/dto"
	"github.com/winstonajr/dashdrop/internal/service"
)

type FileHandler struct {
	svc *service.FileService
}

func NewFileHandler(svc *service.FileService) *FileHandler {
	return &FileHandler{
		svc: svc,
	}
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	files, err := h.svc.ListFiles(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "falha ao buscar arquivos",
		})
		return
	}

	response := make([]dto.FileResponse, 0, len(files))

	for _, f := range files {
		response = append(response, dto.FileResponse{
			ID:   f.ID,
			Name: f.OriginalName,
			Size: f.FileSize,
		})
	}
	c.JSON(http.StatusOK, response)
}

func (h *FileHandler) UploadFiles(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campo 'file' é obrigatório"})
		return
	}

	fileContent, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao ler arquivo"})
		return
	}
	defer fileContent.Close()

	id, err := h.svc.UploadFiles(c.Request.Context(), fileHeader.Filename, fileHeader.Size, fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "upload concluído com sucesso",
		"id":      id,
	})
}

func (h *FileHandler) DownloadFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do arquivo é obrigatório"})
		return
	}
	fullPath, originalName, err := h.svc.GetFileDownload(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "arquivo não encontrado"})
		return
	}
	c.FileAttachment(fullPath, originalName)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id do arquivo é obrigatório"})
		return
	}

	if err := h.svc.DeleteFile(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
