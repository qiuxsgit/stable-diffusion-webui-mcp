package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"runtime/debug"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

func InitMCPServer(appService *AppService) *mcp.Server {
	mcpImpl := &mcp.Implementation{
		Name:    "wechat-article-fetch-mcp",
		Version: "0.0.1",
	}

	server := mcp.NewServer(mcpImpl, nil)

	registerTools(server, appService)

	logrus.Info("MCP Server initialized with official SDK")

	return server
}

func withPanicRecovery[T any](
	toolName string,
	handler func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error),
) func(context.Context, *mcp.CallToolRequest, T) (*mcp.CallToolResult, any, error) {

	return func(ctx context.Context, req *mcp.CallToolRequest, args T) (result *mcp.CallToolResult, resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logrus.WithFields(logrus.Fields{
					"tool":  toolName,
					"panic": r,
				}).Error("Tool handler panicked")

				logrus.Errorf("Stack trace:\n%s", debug.Stack())

				result = &mcp.CallToolResult{
					Content: []mcp.Content{
						&mcp.TextContent{
							Text: fmt.Sprintf("工具 %s 执行时发生内部错误: %v\n\n请查看服务端日志获取详细信息。", toolName, r),
						},
					},
					IsError: true,
				}
				resp = nil
				err = nil
			}
		}()

		return handler(ctx, req, args)
	}
}

// convertToMCPResult 将自定义的 MCPToolResult 转换为官方 SDK 的格式
func convertToMCPResult(result *MCPToolResult) *mcp.CallToolResult {
	var contents []mcp.Content
	for _, c := range result.Content {
		switch c.Type {
		case "text":
			contents = append(contents, &mcp.TextContent{Text: c.Text})
		case "image":
			// 解码 base64 字符串为 []byte
			imageData, err := base64.StdEncoding.DecodeString(c.Data)
			if err != nil {
				logrus.WithError(err).Error("Failed to decode base64 image data")
				// 如果解码失败，添加错误文本
				contents = append(contents, &mcp.TextContent{
					Text: "图片数据解码失败: " + err.Error(),
				})
			} else {
				contents = append(contents, &mcp.ImageContent{
					Data:     imageData,
					MIMEType: c.MimeType,
				})
			}
		}
	}

	return &mcp.CallToolResult{
		Content: contents,
		IsError: result.IsError,
	}
}

func registerTools(mcpServer *mcp.Server, appService *AppService) {
	mcp.AddTool(mcpServer,
		&mcp.Tool{
			Name:        "txt2img",
			Description: "根据文本生成图片",
		},
		withPanicRecovery("text_to_image", func(ctx context.Context, req *mcp.CallToolRequest, arg sdwebui.TextToImageRequest) (*mcp.CallToolResult, any, error) {
			result := appService.mcpHandler.textToImage(ctx, arg)
			return convertToMCPResult(result), nil, nil
		}),
	)

	mcp.AddTool(mcpServer,
		&mcp.Tool{
			Name:        "sd_models",
			Description: "获取SD模型列表",
		},
		withPanicRecovery("sd_models", func(ctx context.Context, req *mcp.CallToolRequest, arg any) (*mcp.CallToolResult, any, error) {
			result := appService.mcpHandler.sdModels(ctx)
			return convertToMCPResult(result), nil, nil
		}),
	)

	mcp.AddTool(mcpServer,
		&mcp.Tool{
			Name:        "switch_model",
			Description: "切换SD模型",
		},
		withPanicRecovery("switch_model", func(ctx context.Context, req *mcp.CallToolRequest, arg sdwebui.SwitchModelRequest) (*mcp.CallToolResult, any, error) {
			result := appService.mcpHandler.switchModel(ctx, arg)
			return convertToMCPResult(result), nil, nil
		}),
	)
}
