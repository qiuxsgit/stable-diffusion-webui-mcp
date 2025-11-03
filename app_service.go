package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/sirupsen/logrus"
	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

type AppService struct {
	sdwebuiService *sdwebui.SdwebuiService
	router         *gin.Engine
	httpServer     *http.Server
	mcpServer      *mcp.Server
	mcpHandler     *McpHandler
	apiHandler     *ApiHandler
}

const BASE_MCP_PATH = "/mcp"

func (s *AppService) Start(port string) error {
	s.router = setupRoutes(s)

	s.httpServer = &http.Server{
		Addr:    port,
		Handler: s.router,
	}

	go func() {
		logrus.Infof("启动 MCP 服务器: http://127.0.0.1%s%s", port, BASE_MCP_PATH)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Errorf("服务器启动失败: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logrus.Infof("正在关闭服务器...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		logrus.Warnf("等待连接关闭超时，强制退出: %v", err)
	} else {
		logrus.Infof("服务器已优雅关闭")
	}

	return nil
}

func setupRoutes(appService *AppService) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())

	router.Use(gin.Recovery())

	setupStreamableHttpHandler(appService, router)

	setupSseEndpoints(appService, router)

	setupApiV1(appService, router)

	return router
}

func setupApiV1(appService *AppService, router *gin.Engine) {
	apiV1Group := router.Group("/api/v1")
	{
		apiV1Group.GET("/read/file/:fileName", appService.apiHandler.readFile)
	}
}

func setupStreamableHttpHandler(appService *AppService, router *gin.Engine) {
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return appService.mcpServer
		},
		&mcp.StreamableHTTPOptions{
			JSONResponse: true,
		},
	)

	router.Any(BASE_MCP_PATH, gin.WrapH(mcpHandler))
	router.Any(fmt.Sprintf("%s/*path", BASE_MCP_PATH), gin.WrapH(mcpHandler))
}

func setupSseEndpoints(appService *AppService, router *gin.Engine) {
	sseMcpHandler := mcp.NewSSEHandler(
		func(request *http.Request) *mcp.Server {
			return appService.mcpServer
		},
		&mcp.SSEOptions{},
	)
	router.Any("/sse", gin.WrapH(sseMcpHandler))
	router.Any("/sse/*path", gin.WrapH(sseMcpHandler))
}

func NewAppService(sdwebuiService *sdwebui.SdwebuiService, apiHandler *ApiHandler) *AppService {
	appService := &AppService{
		sdwebuiService: sdwebuiService,
		mcpHandler:     NewMcpHandler(sdwebuiService),
		apiHandler:     apiHandler,
	}

	appService.mcpServer = InitMCPServer(appService)

	return appService
}
