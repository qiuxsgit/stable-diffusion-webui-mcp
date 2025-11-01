package main

import (
	"flag"

	"qiuxs.com/stable-diffusion-webui-mcp/sdwebui"
)

func main() {
	var (
		port       string
		sdwebuiUrl string
	)

	flag.StringVar(&port, "port", ":18080", "端口")
	flag.StringVar(&sdwebuiUrl, "sdwebuiUrl", "http://127.0.0.1:7860", "Stable Diffusion WebUI 服务地址")

	flag.Parse()

	sdwebuiService := sdwebui.NewSdwebuiService()

}
