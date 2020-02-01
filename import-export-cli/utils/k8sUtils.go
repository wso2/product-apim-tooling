package utils

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// K8sApplyFromFile applies resources from list of files, urls or directories
func K8sApplyFromFile(fileList ...string) error {
	kubectlArgs := []string{K8sApply}
	for _, file := range fileList {
		kubectlArgs = append(kubectlArgs, "-f", file)
	}

	return ExecuteCommand(Kubectl, kubectlArgs...)
}

// K8sApplyFromStdin applies resources from standard input
func K8sApplyFromStdin(stdInput string) error {
	return ExecuteCommandFromStdin(stdInput, Kubectl, K8sApply, "-f", "-")
}

// ExecuteCommand executes the command with args and prints output, errors in standard output, error
func ExecuteCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	setCommandOutAndError(cmd)
	return cmd.Run()
}

// ExecuteCommandFromStdin executes the command with args and prints output the standard output
func ExecuteCommandFromStdin(stdInput string, command string, args ...string) error {
	cmd := exec.Command(command, args...)
	setCommandOutAndError(cmd)

	pipe, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	if _, err = pipe.Write([]byte(stdInput)); err != nil {
		return err
	}
	if err := pipe.Close(); err != nil {
		return err
	}

	return cmd.Run()
}

// GetCommandOutput executes a command and returns the output
func GetCommandOutput(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	setCommandOutAndError(cmd)

	output, err := cmd.Output()
	return string(output), err
}

// setCommandOutAndError sets the output and error of the command cmd to the standard output and error
func setCommandOutAndError(cmd *exec.Cmd) {
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)
}
