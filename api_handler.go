package main

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

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
	// 获取路径参数，去除开头的斜杠
	filePath := c.Param("filePath")
	if len(filePath) > 0 && filePath[0] == '/' {
		filePath = filePath[1:]
	}

	// 安全检查：防止路径遍历攻击
	// 1. 检查是否包含 .. 或 ..\（Windows路径遍历）
	if strings.Contains(filePath, "..") {
		c.String(http.StatusBadRequest, "非法的文件路径")
		return
	}

	// 2. 清理路径并检查是否包含绝对路径或危险字符
	cleanedPath := filepath.Clean(filePath)
	// 如果清理后的路径包含 ..，说明原路径试图越界
	if strings.Contains(cleanedPath, "..") {
		c.String(http.StatusBadRequest, "非法的文件路径")
		return
	}

	// 3. 检查是否以斜杠开头（绝对路径）
	if filepath.IsAbs(cleanedPath) {
		c.String(http.StatusBadRequest, "非法的文件路径")
		return
	}

	// 使用清理后的路径
	file, err := h.fileService.ReadFile(cleanedPath)
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
