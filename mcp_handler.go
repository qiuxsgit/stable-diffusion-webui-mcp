package main

import (
	"context"
	"encoding/json"
	"fmt"

	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

type McpHandler struct {
	sdwebuiService *sdwebui.SdwebuiService
}

func (h *McpHandler) textToImage(ctx context.Context, arg sdwebui.TextToImageRequest) *MCPToolResult {
	response, err := h.sdwebuiService.TextToImage(ctx, arg)
	if err != nil {
		return &MCPToolResult{
			Content: []MCPContent{
				{
					Type: "text",
					Text: "生成图片失败: " + err.Error(),
				},
			},
			IsError: true,
		}
	}

	// 检查是否有生成的图片
	if len(response.Images) == 0 {
		return &MCPToolResult{
			Content: []MCPContent{{
				Type: "text",
				Text: "未生成任何图片",
			}},
			IsError: true,
		}
	}

	// 创建结果内容
	var contents []MCPContent

	contents = append(contents, MCPContent{
		Type: "text",
		Text: fmt.Sprintf("图片生成成功！数量: %d", len(response.Images)),
	})

	if response.Parameters != nil {
		paramJson, _ := json.Marshal(response.Parameters)
		contents = append(contents, MCPContent{
			Type: "text",
			Text: fmt.Sprintf("生成参数: %s", paramJson),
		})
	}

	// 添加生成信息
	if response.Info != "" {
		contents = append(contents, MCPContent{
			Type: "text",
			Text: fmt.Sprintf("扩展信息: %s", response.Info),
		})
	}

	// 添加生成的图片
	for _, imageUrl := range response.Images {
		contents = append(contents, MCPContent{
			Type: "text",
			Text: imageUrl,
		})
	}

	return &MCPToolResult{
		Content: contents,
	}
}

func (h *McpHandler) sdModels(ctx context.Context) *MCPToolResult {
	models, err := h.sdwebuiService.SdModels(ctx)
	if err != nil {
		return errorResult(fmt.Sprintf("获取模型列表失败: %v", err))
	}
	jsonModels, err := json.Marshal(models)
	if err != nil {
		return errorResult(fmt.Sprintf("获取模型列表失败: %v", err))
	}
	return successResult(toContents(makeTextContent(string(jsonModels))))
}

func (h *McpHandler) switchModel(ctx context.Context, arg sdwebui.SwitchModelRequest) *MCPToolResult {
	response, err := h.sdwebuiService.SwitchModel(ctx, arg)
	if err != nil {
		return errorResult(fmt.Sprintf("切换模型失败: %v", err))
	}

	if response.Success {
		return successResult(toContents(makeTextContent(response.Message)))
	}

	return errorResult(response.Message)
}

func toContents(content ...MCPContent) []MCPContent {
	return content
}

func makeTextContent(text string) MCPContent {
	return MCPContent{
		Type: "text",
		Text: text,
	}
}

func successResult(contents []MCPContent) *MCPToolResult {
	return &MCPToolResult{
		Content: contents,
		IsError: false,
	}
}

func errorResult(message string) *MCPToolResult {
	return &MCPToolResult{
		Content: []MCPContent{
			{
				Type: "text",
				Text: message,
			},
		},
		IsError: true,
	}
}

func NewMcpHandler(sdwebuiService *sdwebui.SdwebuiService) *McpHandler {
	return &McpHandler{
		sdwebuiService: sdwebuiService,
	}
}
