package sdwebui

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"qiuxs.com/stable-diffusion-webui-mcp/internal"
)

type SdwebuiService struct {
	baseUrl     string
	fileService *internal.FileService
	client      *http.Client
}

func NewSdwebuiService(sdwebuiUrl string, fileService *internal.FileService) *SdwebuiService {
	return &SdwebuiService{
		baseUrl:     sdwebuiUrl,
		fileService: fileService,
		client: &http.Client{
			Timeout: 300 * time.Second, // 5分钟超时，因为图片生成可能需要较长时间
		},
	}
}

func (s *SdwebuiService) TextToImage(ctx context.Context, arg TextToImageRequest) (*TextToImageResponse, error) {
	// 设置默认值
	if arg.Width == 0 {
		arg.Width = 512
	}
	if arg.Height == 0 {
		arg.Height = 512
	}
	if arg.Steps == 0 {
		arg.Steps = 20
	}
	if arg.SamplerName == "" {
		arg.SamplerName = "Euler a"
	}
	if arg.CFGScale == 0 {
		arg.CFGScale = 7.0
	}
	if arg.BatchSize == 0 {
		arg.BatchSize = 1
	}
	if arg.NIter == 0 {
		arg.NIter = 1
	}

	// 准备请求体
	requestBody, err := json.Marshal(arg)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %v", err)
	}

	// 构建API URL
	apiUrl := fmt.Sprintf("%s/sdapi/v1/txt2img", s.baseUrl)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用Stable Diffusion API失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	var response TextToImageResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %v", err)
	}

	// 保存生成的图片
	var fileUrls []string
	for _, imageData := range response.Images {
		fileUrl, err := s.fileService.SaveImage(imageData)
		if err != nil {
			return nil, fmt.Errorf("保存图片失败: %v", err)
		}
		fileUrls = append(fileUrls, fileUrl)
	}

	// 将保存的图片路径添加到响应中
	response.Images = fileUrls

	return &response, nil
}

func (s *SdwebuiService) SdModels(ctx context.Context) (*SdModelsResponse, error) {
	// 构建API URL
	apiUrl := fmt.Sprintf("%s/sdapi/v1/sd-models", s.baseUrl)

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", apiUrl, nil)
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用Stable Diffusion API失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	// 解析响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	var models []SdModel

	if err := json.Unmarshal(body, &models); err != nil {
		return nil, fmt.Errorf("解析响应JSON失败: %v", err)
	}

	return &SdModelsResponse{
		Models: models,
	}, nil
}

func (s *SdwebuiService) SwitchModel(ctx context.Context, arg SwitchModelRequest) (*SwitchModelResponse, error) {
	// 构建API URL
	apiUrl := fmt.Sprintf("%s/sdapi/v1/options", s.baseUrl)

	// 准备请求参数
	requestBody, err := json.Marshal(arg)
	if err != nil {
		return nil, fmt.Errorf("序列化请求参数失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "POST", apiUrl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("调用Stable Diffusion API失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API返回错误状态码: %d %s, 响应: %s", resp.StatusCode, resp.Status, body)
	}

	return &SwitchModelResponse{
		Success: true,
		Message: "切换模型成功",
	}, nil
}
