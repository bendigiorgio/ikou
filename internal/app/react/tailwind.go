package react

import (
	"fmt"
	"os/exec"
)

func BuildCSS(cssPath string, outPath string) error {
	cmdTest := exec.Command("pwd")
	out, err := cmdTest.Output()
	if err != nil {
		return fmt.Errorf("error getting current working directory: %w", err)
	}
	fmt.Println(string(out))
	cmd := exec.Command("./tailwindcss", "-i", cssPath, "-o", outPath, "-c", "../../../frontend/tailwind.config.js")
	cmd.Dir = "./internal/app/lib"
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error building css: %w", err)
	}
	return nil
}
