package server

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func RunDeepSeek(context string) (string, error) {
	//Log the context being passed to the model
	fmt.Println("🧾 Context being passed to DeepSeek:\n--------------------------------")
	fmt.Println(context)
	fmt.Println("--------------------------------")

	// Write context to a temp file
	tmpFile, err := os.CreateTemp("", "deepseek_input_*.txt")
	if err != nil {
		return "", fmt.Errorf("error creating temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(context)
	tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("error writing to temp file: %w", err)
	}

	// Run DeepSeek model
	cmd := exec.Command("C:\\Users\\aikaw\\AppData\\Local\\Programs\\Ollama\\ollama.exe", "run", "deepseek-r1:8b")
	cmd.Stdin = strings.NewReader(context)

	var outBuffer, errBuffer bytes.Buffer
	cmd.Stdout = &outBuffer
	cmd.Stderr = &errBuffer

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("DeepSeek error: %s\n%s", err, errBuffer.String())
	}

	return CleanANSI(outBuffer.String()), nil
}
