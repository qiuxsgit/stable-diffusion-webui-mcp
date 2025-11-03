package internal

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	filePath := filepath.Join(s.fileSavePath, fileName)

	// 确保目录存在
	if err := os.MkdirAll(s.fileSavePath, 0755); err != nil {
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

	fileUrl := fmt.Sprintf("%s/api/v1/read/file/%s", s.serverUrl, fileName)
	logrus.Infof("fileUrl: %s", fileUrl)

	return fileUrl, nil
}

func (s *FileService) ReadFile(fileName string) (*os.File, error) {
	if strings.Contains(fileName, "..") {
		return nil, errors.New("file name contains invalid characters")
	}
	filePath := s.fileSavePath + "/" + fileName
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}
