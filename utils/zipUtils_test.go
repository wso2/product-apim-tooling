package utils

import "testing"

func TestZipDir(t *testing.T) {
	err := ZipDir("", "")
	if err == nil {
		t.Errorf("ZipDir() didn't return an error for invalid source and destination")
	}
}