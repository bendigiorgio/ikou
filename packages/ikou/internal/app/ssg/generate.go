package ssg

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/router"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"go.uber.org/zap"
)

func GenerateStaticSite() error {
	outputDir := utils.GlobalConfig.OutPath

	basePath := utils.GlobalConfig.BasePath
	staticPath := utils.GlobalConfig.StaticPath
	useSrc := utils.GlobalConfig.UseSrc

	srcPath := basePath
	if useSrc {
		srcPath = path.Join(basePath, "src")
	}

	router.InitializeRouting(srcPath, false)

	if err := utils.CopyDir(staticPath, filepath.Join(outputDir, "public")); err != nil {
		utils.Logger.Error("Error copying static files", zap.Error(err))
		return err
	}

	for route, routeInfo := range router.RouteMap {
		initialProps := react.PageProps{
			PageRoute: route,
		}

		pageData, err := react.RenderPage(routeInfo.IsSSG, initialProps, routeInfo.PagePath)
		if err != nil {
			utils.Logger.Error("Error rendering page", zap.String("route", route), zap.Error(err))
			return err
		}

		outputPath := filepath.Join(outputDir, route)
		if route == "/" || route == "." {
			outputPath = filepath.Join(outputDir, "index.html")
		} else {
			if err := os.MkdirAll(filepath.Dir(outputPath), fs.ModePerm); err != nil {
				utils.Logger.Error("Error creating directories", zap.String("path", outputPath), zap.Error(err))
				return err
			}
			if filepath.Ext(outputPath) == "" {
				outputPath += ".html"
			}
		}

		file, err := os.Create(outputPath)
		if err != nil {
			utils.Logger.Error("Error creating file", zap.String("path", outputPath), zap.Error(err))
			return err
		}
		defer file.Close()

		err = pageData.Tmpl.Execute(file, pageData)
		if err != nil {
			utils.Logger.Error("Error writing template to file", zap.String("path", outputPath), zap.Error(err))
			return err
		}

		utils.Logger.Info("Generated static page", zap.String("path", outputPath))
	}

	return nil
}
