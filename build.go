// build.go
package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func main() {
	cmdDir := "./cmd/gitinsight"
	// 输出目录
	outputDir := "bin"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		return
	}

	// 日期版本号，例如 20251007
	dateTag := time.Now().Format("20060102")

	targets := []struct {
		GOOS   string
		GOARCH string
	}{
		{"darwin", "amd64"},
		{"linux", "amd64"},
		{"windows", "amd64"},
	}

	for _, t := range targets {
		outputName := fmt.Sprintf("gitinsight-%s-%s-%s", t.GOOS, t.GOARCH, dateTag)
		if t.GOOS == "windows" {
			outputName += ".exe"
		}
		outputPath := filepath.Join(outputDir, outputName)

		fmt.Printf("Building for %s/%s -> %s\n", t.GOOS, t.GOARCH, outputPath)

		cmd := exec.Command("go", "build", "-o", outputPath, cmdDir)
		cmd.Env = append(os.Environ(),
			"CGO_ENABLED=0",
			"GOOS="+t.GOOS,
			"GOARCH="+t.GOARCH,
		)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Build failed for %s/%s: %v\n", t.GOOS, t.GOARCH, err)
			continue
		}
		fmt.Printf("Build succeeded: %s\n", outputPath)
	}
}
