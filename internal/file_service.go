package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FileService struct {
	fileSavePath string
	serverUrl    string
}

func NewFileService(fileSavePath string, serverUrl string) *FileService {
	return &FileService{
		fileSavePath: fileSavePath,
		serverUrl:    serverUrl,
	}
}

// saveImage 将base64图片数据保存到指定路径
func (s *FileService) SaveImage(base64Data string) (string, error) {
	// 生成UUID作为文件名
	fileID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("生成UUID失败: %v", err)
	}

	// 创建文件名
	fileName := fmt.Sprintf("%s.png", fileID.String())
	// 生成日期文件夹（yyyy-MM-dd格式）
	dateFolder := time.Now().Format("2006-01-02")
	// 构建完整路径，包含日期文件夹
	filePath := filepath.Join(s.fileSavePath, dateFolder, fileName)

	// 确保目录存在（包括日期文件夹）
	dateDir := filepath.Join(s.fileSavePath, dateFolder)
	if err := os.MkdirAll(dateDir, 0755); err != nil {
		return "", fmt.Errorf("创建目录失败: %v", err)
	}

	// 解码base64数据
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return "", fmt.Errorf("解码base64数据失败: %v", err)
	}

	// 保存图片文件
	if err := os.WriteFile(filePath, imageData, 0644); err != nil {
		return "", fmt.Errorf("保存图片文件失败: %v", err)
	}

	// 构建相对路径用于URL（日期文件夹/文件名），使用path包确保URL使用正斜杠
	relativePath := path.Join(dateFolder, fileName)
	fileUrl := fmt.Sprintf("%s/api/v1/read/file/%s", s.serverUrl, relativePath)
	logrus.Infof("fileUrl: %s", fileUrl)

	return fileUrl, nil
}

func (s *FileService) ReadFile(filePath string) (*os.File, error) {
	if strings.Contains(filePath, "..") {
		return nil, errors.New("file path contains invalid characters")
	}
	// 使用filepath.Join安全地拼接路径，支持日期文件夹/文件名的格式
	fullPath := filepath.Join(s.fileSavePath, filePath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
