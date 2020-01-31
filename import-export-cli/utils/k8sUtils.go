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

	cmd := exec.Command("kubectl", kubectlArgs...)
	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

	return cmd.Run()
}

// K8sApplyFromStdin applies resources from standard input
func K8sApplyFromStdin(stdInput string) error {
	cmd := exec.Command(
		Kubectl,
		K8sApply,
		"-f",
		"-",
	)

	var errBuf, outBuf bytes.Buffer
	cmd.Stderr = io.MultiWriter(os.Stderr, &errBuf)
	cmd.Stdout = io.MultiWriter(os.Stdout, &outBuf)

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
