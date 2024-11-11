package router

import (
	"net/http"
	"os"
	"path/filepath"
	"plugin"

	"github.com/bendigiorgio/ikou/internal/app/utils"
	"github.com/fsnotify/fsnotify"
)

const MIDDLEWARE_PATH = "routes/middleware/middleware.go"

var GLOBAL_MIDDLEWARE func(http.Handler) http.Handler

func loadMiddleware() error {
	if _, err := os.Stat(MIDDLEWARE_PATH); os.IsNotExist(err) {
		return nil
	}
	pluginPath, err := compileToPlugin(MIDDLEWARE_PATH)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to compile middleware to plugin: %v", err)
		return err
	}

	p, err := plugin.Open(pluginPath)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to load middleware plugin: %v", err)
		return err
	}

	middlewareSymbol, err := p.Lookup("Middleware")
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to find Middleware function in plugin: %v", err)
		return err
	}

	middlewareFunc, ok := middlewareSymbol.(func(http.Handler) http.Handler)
	if !ok {
		utils.Logger.Sugar().Errorf("Middleware function has an incorrect signature")
		return err
	}

	GLOBAL_MIDDLEWARE = middlewareFunc
	utils.Logger.Sugar().Info("Loaded middleware successfully")
	return nil
}

func watchMiddlewareDirectory() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to create middleware directory watcher: %v", err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					utils.Logger.Sugar().Infof("Middleware file changed: %s", event.Name)
					err := loadMiddleware()
					if err != nil {
						utils.Logger.Sugar().Errorf("Failed to reload middleware: %v", err)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Logger.Sugar().Errorf("Middleware directory watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(filepath.Dir(MIDDLEWARE_PATH))
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to watch middleware directory: %v", err)
	}
	select {}
}
