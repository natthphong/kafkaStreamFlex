package test

import (
	"context"
	"fmt"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/internal/script"
	"gitlab.com/home-server7795544/home-server/iam/iam-backend/sdk"
	"os"
	"path/filepath"
	"testing"
)

func TestExecuteBasic(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get current directory:", err)
	}
	parentDir := filepath.Dir(cwd)

	sourcePath := filepath.Join(parentDir, "sdk/test/example_script.go")
	scriptPath := filepath.Join(parentDir, "sdk/test/example_script.so")

	err = script.BuildPlugin(sourcePath, scriptPath)
	if err != nil {
		return
	}
	payload := []byte(`{"hello":"world"}`)

	client := sdk.NewStreamFlexClient(
		context.Background(),
		nil,
		nil,
		nil,
		nil,
		nil,
		payload,
	)

	err = script.ExecuteScript(scriptPath, client)
	if err != nil {
		fmt.Println("Script execution failed:", err)
	} else {
		fmt.Println("Script executed successfully!")
	}
}
