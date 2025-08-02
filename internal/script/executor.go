package script

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"

	sdk "github.com/natthphong/streamFlexSdk/client"
)

func BuildPlugin(src, dst string) error {
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}

	cmd := exec.Command(
		"go", "build",
		"-buildmode=plugin",
		"-o", dst,
		src,
	)

	cmd.Env = append(os.Environ(),
		"CGO_ENABLED=1",
		"GOTOOLCHAIN=local",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("build plugin: %w\n%s", err, out)
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
