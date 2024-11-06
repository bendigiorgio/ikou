package router

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
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

type ApiHandlerFn func(http.ResponseWriter, *http.Request, string)

type ApiRouteInfo struct {
	FilePath  string
	Method    string
	HandlerFn ApiHandlerFn
}

type EntryRouteFn func(http.ResponseWriter, *http.Request, string) map[string]interface{}

type EntryRouteInfo struct {
	FilePath  string
	HandlerFn EntryRouteFn
	Route     *RouteInfo
}

var RouteMap = map[string]RouteInfo{}
var ApiRouteMap = map[string]ApiRouteInfo{}
var EntryRouteMap = map[string]EntryRouteInfo{}

const BASE_API_ROUTE = "routes/api"
const BASE_ENTRY_ROUTE = "routes/entry"

// Automatically compiles a .go file to a .so file
func compileToPlugin(filePath string) (string, error) {
	outputPath := strings.TrimSuffix(filePath, ".go") + ".so"
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", outputPath, filePath)
	err := cmd.Run()
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to compile %s to plugin: %v", filePath, err)
		return "", err
	}
	utils.Logger.Sugar().Infof("Compiled %s to %s", filePath, outputPath)
	return outputPath, nil
}

// Scan and generate routes for API handlers
func scanApiDirectory() error {
	return filepath.Walk(BASE_API_ROUTE, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Compile and generate API route
			pluginPath, err := compileToPlugin(path)
			if err == nil {
				generateApiRoute(pluginPath)
			}
		}
		return nil
	})
}

func generateApiRoute(filePath string) {
	apiPath := utils.GlobalConfig.ApiPath

	fileName := filepath.Base(filePath)
	method := strings.ToUpper(strings.TrimSuffix(fileName, filepath.Ext(fileName))) // e.g., "get" becomes "GET"

	// Generate route by removing BASE_API_ROUTE and the method part of the path
	routeDir := filepath.Dir(strings.TrimPrefix(filePath, BASE_API_ROUTE+"/"))
	route := apiPath + "/" + filepath.ToSlash(routeDir)

	// Load the plugin and look up the Handler function
	p, err := plugin.Open(filePath)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to load API plugin %s: %v", filePath, err)
		return
	}
	handlerSymbol, err := p.Lookup("Handler")
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to find Handler in %s: %v", filePath, err)
		return
	}
	handler, ok := handlerSymbol.(func(http.ResponseWriter, *http.Request, string))
	if !ok {
		utils.Logger.Sugar().Errorf("Handler in %s has an incorrect signature", filePath)
		return
	}

	ApiRouteMap[route] = ApiRouteInfo{
		FilePath:  filePath,
		Method:    method,
		HandlerFn: handler,
	}

	utils.Logger.Sugar().Debugf("Mapped API route: %s %s -> %s", method, route, filePath)
}

// Scan and generate routes for Entry handlers
func scanEntryDirectory() error {
	return filepath.Walk(BASE_ENTRY_ROUTE, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".go") {
			// Compile and generate Entry route
			pluginPath, err := compileToPlugin(path)
			if err == nil {
				generateEntryRoute(pluginPath)
			}
		}
		return nil
	})
}

func generateEntryRoute(filePath string) {
	route := strings.TrimSuffix(strings.TrimPrefix(filePath, BASE_ENTRY_ROUTE+"/"), ".so")

	// Load the plugin and look up the Entry function
	p, err := plugin.Open(filePath)
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to load Entry plugin %s: %v", filePath, err)
		return
	}
	entrySymbol, err := p.Lookup("Entry")
	if err != nil {
		utils.Logger.Sugar().Errorf("Failed to find Entry in %s: %v", filePath, err)
		return
	}
	entryHandler, ok := entrySymbol.(func(http.ResponseWriter, *http.Request, string) map[string]interface{})
	if !ok {
		utils.Logger.Sugar().Errorf("Entry in %s has an incorrect signature", filePath)
		return
	}

	if pageRoute, exists := RouteMap[route]; exists {
		EntryRouteMap[route] = EntryRouteInfo{
			FilePath:  filePath,
			HandlerFn: entryHandler,
			Route:     &pageRoute,
		}
		utils.Logger.Sugar().Debugf("Mapped entry route: %s -> %s", route, filePath)
	} else {
		utils.Logger.Sugar().Warnf("Entry route %s has no matching page route", route)
	}
}

// Watch API directory for changes and recompile as needed
func watchApiDirectory() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to create API directory watcher: %v", err)
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
					utils.Logger.Sugar().Infof("API file changed: %s", event.Name)
					pluginPath, err := compileToPlugin(event.Name)
					if err == nil {
						generateApiRoute(pluginPath)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Logger.Sugar().Errorf("API directory watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(BASE_API_ROUTE)
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to watch API directory: %v", err)
	}
	select {}
}

// Watch Entry directory for changes and recompile as needed
func watchEntryDirectory() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to create Entry directory watcher: %v", err)
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
					utils.Logger.Sugar().Infof("Entry file changed: %s", event.Name)
					pluginPath, err := compileToPlugin(event.Name)
					if err == nil {
						generateEntryRoute(pluginPath)
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				utils.Logger.Sugar().Errorf("Entry directory watcher error: %v", err)
			}
		}
	}()

	err = watcher.Add(BASE_ENTRY_ROUTE)
	if err != nil {
		utils.Logger.Sugar().Fatalf("Failed to watch Entry directory: %v", err)
	}
	select {}
}

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

	err = scanApiDirectory()
	if err != nil {
		utils.Logger.Sugar().Fatalf("Error scanning API directory: %v", err)
	}

	err = scanEntryDirectory()
	if err != nil {
		utils.Logger.Sugar().Fatalf("Error scanning Entry directory: %v", err)
	}

	if dev {
		go watchDirectory(fmt.Sprintf("%s/pages/", baseRoute), baseRoute)
		go watchApiDirectory()
		go watchEntryDirectory()
	}

	utils.Logger.Sugar().Debugf("Initial routes:", RouteMap)
}
