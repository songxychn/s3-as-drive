package main

import (
	"context"
	"embed"
	"s3-as-drive/backend/service"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	configService := service.NewConfigService()
	fileService := service.NewFileService()

	err := wails.Run(&options.App{
		Title:  "s3-as-drive",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			configService.Startup(ctx)
			fileService.Startup(ctx)
		},
		Bind: []interface{}{
			configService,
			fileService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
