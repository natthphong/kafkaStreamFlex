package script

import (
	"fmt"
	sdk "github.com/natthphong/streamFlexSdk/client"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
)

func BuildPlugin(sourceFilePath, outputSoPath string) error {
	outDir := filepath.Dir(outputSoPath)
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	cmd := exec.Command(
		"go",
		"build",
		"-buildmode=plugin",
		"-o", outputSoPath,
		sourceFilePath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("build plugin failed: %w\nOutput:\n%s", err, string(output))
	}

	return nil
}

func ExecuteScript(scriptPath string, client *sdk.StreamFlexClient) error {
	plg, err := plugin.Open(scriptPath)
	if err != nil {
		// TODO not found go load from s3 server
		return fmt.Errorf("failed to load plugin %s: %w", scriptPath, err)
	}

	symProcess, err := plg.Lookup("Process")
	if err != nil {
		return fmt.Errorf("function 'Process' not found in plugin: %w", err)
	}
	processFunc, ok := symProcess.(func(*sdk.StreamFlexClient) error)
	if !ok {
		return fmt.Errorf("'Process' has incompatible type in plugin %s", scriptPath)
	}

	err = processFunc(client)
	if err != nil {
		log.Printf("Error during script execution: %v\n", err)
		return err
	}

	log.Printf("Script %s executed successfully!", scriptPath)
	return nil
}
