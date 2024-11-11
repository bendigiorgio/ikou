package router

import (
	"os/exec"
	"strings"

	"github.com/bendigiorgio/ikou/internal/app/utils"
)

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
