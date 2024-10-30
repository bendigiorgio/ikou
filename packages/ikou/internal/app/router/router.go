package router

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bendigiorgio/ikou/internal/app/react"
	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/fsnotify/fsnotify"
)

type RouteInfo struct {
	PagePath     string
	IsSSG        bool
	IsDynamic    bool
	DynamicNames []string
}

var RouteMap = map[string]RouteInfo{}

func scanDirectory(directory string, baseRoute string) error {
	return filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(path, ".page.tsx") || strings.HasSuffix(path, ".page.jsx")) {
			route := generateRouteFromFilePath(path, baseRoute)

			isSSG := !strings.Contains(path, ".client.page")

			if !isSSG {
				route = strings.Replace(route, ".client", "", -1)
			}

			RouteMap[route] = RouteInfo{
				PagePath: path,
				IsSSG:    isSSG,
			}

			utils.Logger.Sugar().Debugf("Mapped route: %s -> %s (SSR: %v)\n", route, path, isSSG)
		}
		return nil
	})
}

func generateRouteFromFilePath(filePath string, baseRoute string) string {
	route := strings.TrimPrefix(filePath, baseRoute)
	route = strings.TrimPrefix(route, "/pages/")
	route = strings.TrimSuffix(route, ".page.tsx")
	route = strings.TrimSuffix(route, ".page.jsx")

	if route == "index" {
		return "/"
	}

	route = strings.Replace(route, "index", "/", -1)

	return "/" + route
}

func watchDirectory(directory string, baseRoute string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					utils.Logger.Sugar().Infof("File changed: ", event.Name)
					err := scanDirectory(directory, baseRoute)
					if err != nil {
						utils.Logger.Sugar().Errorf("Error rescanning directory:", err)
					}
					err = react.BuildCSS()
					if err != nil {
						utils.Logger.Sugar().Fatalf("Error building CSS: %v", err)

					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Logger.Sugar().Errorf("Watcher error:", err)
			}
		}
	}()

	err = watcher.Add(directory)
	if err != nil {
		utils.Logger.Sugar().Fatal(err)
	}
	<-done
}

func InitializeRouting(baseRoute string, dev bool) {
	err := scanDirectory(fmt.Sprintf("%s/pages/", baseRoute), baseRoute)
	if err != nil {
		utils.Logger.Sugar().Fatalf("Error scanning pages directory: %v", err)
	}
	if dev {
		go watchDirectory(fmt.Sprintf("%s/pages/", baseRoute), baseRoute)
	}

	utils.Logger.Sugar().Debugf("Initial routes:", RouteMap)
}
