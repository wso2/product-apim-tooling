package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIsValid1 - Create new file
func TestIsValid1(t *testing.T) {
	filePath := filepath.Join(CurrentDir, "test.txt")
	IsValid(filePath)
}

// TestIsValid2 - Create new file
func TestIsValid2(t *testing.T) {
	fileName := "test.txt"
	os.Create(fileName)
	filePath := filepath.Join(CurrentDir, fileName)
	IsValid(filePath)
	os.Remove(filePath)
}
