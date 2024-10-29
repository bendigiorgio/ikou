package react

import (
	"fmt"
	"os/exec"
	"path"

	"github.com/bendigiorgio/ikou/internal/app/utils"
)

func BuildCSS() error {
	defer utils.Logger.Sync()

	cssPath := utils.GlobalConfig.Tailwind.CSSPath
	outPath := utils.GlobalConfig.Tailwind.Output
	configPath := utils.GlobalConfig.Tailwind.Config
	basePath := utils.GlobalConfig.BasePath

	cssPath = path.Join(basePath, cssPath)
	outPath = path.Join(basePath, outPath)
	configPath = path.Join(basePath, configPath)

	tailwindExecutable := "./internal/app/lib/tailwindcss"

	cmd := exec.Command(tailwindExecutable, "-i", cssPath, "-o", outPath, "-c", configPath)

	output, err := cmd.CombinedOutput()

	utils.Logger.Sugar().Debugf("Tailwind CSS Build Output:\n%s", output)

	if err != nil {
		return fmt.Errorf("error building css: %w", err)
	}
	return nil
}

func WatchCSS() error {
	defer utils.Logger.Sync()

	cssPath := utils.GlobalConfig.Tailwind.CSSPath
	outPath := utils.GlobalConfig.Tailwind.Output
	configPath := utils.GlobalConfig.Tailwind.Config
	basePath := utils.GlobalConfig.BasePath

	cssPath = path.Join(basePath, cssPath)
	outPath = path.Join(basePath, outPath)
	configPath = path.Join(basePath, configPath)

	tailwindExecutable := "./internal/app/lib/tailwindcss"

	cmd := exec.Command(tailwindExecutable, "-i", cssPath, "-o", outPath, "-c", configPath, "--watch")

	output, err := cmd.CombinedOutput()

	utils.Logger.Sugar().Debugf("Tailwind CSS Build Output:\n%s", output)

	if err != nil {
		return fmt.Errorf("error building css: %w", err)
	}
	return nil
}
