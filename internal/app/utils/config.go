package utils

import (
	"encoding/json"
	"os"

	"github.com/fsnotify/fsnotify"
)

var GlobalConfig IkouConfig

type IkouConfig struct {
	BasePath   string `json:"basePath"`
	OutPath    string `json:"outputPath"`
	StaticPath string `json:"staticPath"`
	UseSrc     bool   `json:"useSrc"`
	Port       int    `json:"port"`
	Tailwind   struct {
		Config  string `json:"config"`
		CSSPath string `json:"cssPath"`
		Output  string `json:"output"`
	} `json:"tailwind"`
	LogPath string `json:"logPath"`
}

const BaseJSONConfig = `{
  "basePath": "./frontend",
  "outputPath": "./dist",
  "staticPath": "frontend/public",
  "useSrc": true,
  "port": 3000,
  "tailwind": {
    "config": "tailwind.config.js",
    "cssPath": "src/styles/base.css",
    "output": "public/style.css"
  },
  "logPath": "storage/logs/ikou.log"
}`

func ExtractConfigDetails(configPath string) {
	defer Logger.Sync()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := json.Unmarshal([]byte(BaseJSONConfig), &GlobalConfig); err != nil {
			Logger.Sugar().Fatalf("failed to unmarshal base JSON config: %v", err)
		}
		return
	}

	file, err := os.Open(configPath)
	if err != nil {
		Logger.Sugar().Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config IkouConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		Logger.Sugar().Fatalf("failed to decode config file: %v", err)
	}

	GlobalConfig = config

	Logger.Sugar().Debug("Config loaded: %+v\n", GlobalConfig)
}

func WatchForConfigChanges(configPath string) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Logger.Sugar().Fatalf("error creating file watcher: %v", err)
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
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename {
					Logger.Sugar().Infof("Config file changed, reloading...")
					ExtractConfigDetails(configPath)
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				Logger.Sugar().Errorf("error watching config file: %v", err)
			}
		}
	}()

	err = watcher.Add(configPath)
	if err != nil {
		Logger.Sugar().Fatalf("error adding watcher to config file: %v", err)
	}

	<-done
}
