package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"qiuxs.com/stable-diffusion-webui-mcp/internal"
	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

func main() {
	var (
		port          string
		sdwebuiUrl    string
		imageSavePath string
		serverUrl     string
	)

	flag.StringVar(&port, "port", ":18080", "端口")
	flag.StringVar(&sdwebuiUrl, "sdwebui-url", "http://127.0.0.1:7860", "Stable Diffusion WebUI 服务地址")
	flag.StringVar(&imageSavePath, "image-save-path", "./images", "生成的图片存储位置")
	flag.StringVar(&serverUrl, "server-url", "http://127.0.0.1:18080/", "访问MCP服务的url")

	flag.Parse()

	logrus.Infof("using Stable Diffusion WebUI server: %s", sdwebuiUrl)
	logrus.Infof("save image to: %s", imageSavePath)
	logrus.Infof("server url: %s", serverUrl)

	fileService := internal.NewFileService(imageSavePath, serverUrl)

	sdwebuiService := sdwebui.NewSdwebuiService(sdwebuiUrl, fileService)

	apiHandler := NewApiHandler(fileService)
	appService := NewAppService(sdwebuiService, apiHandler)

	if err := appService.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
