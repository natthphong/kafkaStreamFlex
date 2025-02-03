package test

import (
	"context"
	"fmt"
	"github.com/natthphong/kafkaStreamFlex/internal/script"
	"github.com/natthphong/streamFlexSdk/client"
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteBasic(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get current directory:", err)
	}

	sourcePath := filepath.Join(cwd, "sdk/example_script.go")

	fmt.Println("sourcePath:", sourcePath)
	scriptPath := filepath.Join(cwd, "sdk/example_script.so")

	err = script.BuildPlugin(sourcePath, scriptPath)
	if err != nil {
		t.Fatal("Failed to build plugin:", err)
		return
	}
	payload := []byte(`{"hello":"world"}`)
	fmt.Println("scriptPath:", scriptPath)
	streamClient := client.NewStreamFlexClient(
		context.Background(),
		nil,
		nil,
		nil,
		nil,
		nil,
		payload,
	)

	err = script.ExecuteScript(scriptPath, streamClient)
	if err != nil {
		fmt.Println("Script execution failed:", err)
	} else {
		fmt.Println("Script executed successfully!")
	}
}
