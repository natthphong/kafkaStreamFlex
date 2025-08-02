package test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/natthphong/kafkaStreamFlex/internal/script"
	"github.com/natthphong/streamFlexSdk/client"
)

func TestExecuteBasic(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal("Failed to get current directory:", err)
	}

	sourcePath := filepath.Join(cwd, "sdk/")

	fmt.Println("sourcePath:", sourcePath)
	scriptPath := filepath.Join(cwd, "sdk/example_script.so")

	//err = script.BuildPlugin(sourcePath, scriptPath)
	//if err != nil {
	//	t.Fatal("Failed to build plugin:", err)
	//	return
	//}
	payload := []byte(`{"hello":"world"}`)
	fmt.Println("scriptPath:", scriptPath)
	streamClient := client.NewStreamFlexClient(
		context.Background(),
		nil,
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
