package main

import (
	"flag"

	"github.com/sirupsen/logrus"
	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

func main() {
	var (
		port       string
		sdwebuiUrl string
	)

	flag.StringVar(&port, "port", ":18080", "端口")
	flag.StringVar(&sdwebuiUrl, "sdwebui-url", "http://127.0.0.1:7860", "Stable Diffusion WebUI 服务地址")

	flag.Parse()

	logrus.Infof("using Stable Diffusion WebUI server: %s", sdwebuiUrl)

	sdwebuiService := sdwebui.NewSdwebuiService(sdwebuiUrl)
	appService := NewAppService(sdwebuiService)

	if err := appService.Start(port); err != nil {
		logrus.Fatalf("failed to run server: %v", err)
	}
}
