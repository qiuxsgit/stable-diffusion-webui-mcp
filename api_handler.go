package main

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"qiuxs.com/stable-diffusion-webui-mcp/internal"
)

type ApiHandler struct {
	fileService *internal.FileService
}

func NewApiHandler(fileService *internal.FileService) *ApiHandler {
	return &ApiHandler{
		fileService: fileService,
	}
}

func (h *ApiHandler) readFile(c *gin.Context) {
	file, err := h.fileService.ReadFile(c.Param("fileName"))
	if err != nil {
		c.String(http.StatusNotFound, "文件不存在")
		return
	}

	fileData, err := io.ReadAll(file)
	if err != nil {
		c.String(http.StatusInternalServerError, "读取文件失败")
		return
	}

	c.Data(http.StatusOK, "image/png", fileData)
}
