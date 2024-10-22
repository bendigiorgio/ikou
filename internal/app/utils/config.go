package utils

import (
	"encoding/json"
	"log"
	"os"
)

var GlobalConfig IkouConfig

type IkouConfig struct {
	BasePath string `json:"basePath"`
	OutPath  string `json:"outputPath"`
	UseSrc   bool   `json:"useSrc"`
	Tailwind struct {
		Config  string `json:"config"`
		CSSPath string `json:"cssPath"`
		Output  string `json:"output"`
	} `json:"tailwind"`
}

const BaseJSONConfig = `{
  "basePath": "./frontend",
  "outputPath": "./dist",
  "useSrc": true,
  "tailwind": {
    "config": "tailwind.config.js",
    "cssPath": "src/styles/base.css",
    "output": "public/style.css"
  }
}`

func ExtractConfigDetails(configPath string) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if err := json.Unmarshal([]byte(BaseJSONConfig), &GlobalConfig); err != nil {
			log.Fatalf("failed to unmarshal base JSON config: %v", err)
		}
		return
	}

	file, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config IkouConfig
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		log.Fatalf("failed to decode config file: %v", err)
	}

	GlobalConfig = config
}
