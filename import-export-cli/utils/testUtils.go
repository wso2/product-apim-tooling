package utils

import "os"

// GetRelativeTestDataPathFromImpl
// returns the relative path of the cmd/testdata folder (where artifacts used for testing reside)
func GetRelativeTestDataPathFromImpl() string {
    return ".." + string(os.PathSeparator) + "cmd" + string(os.PathSeparator) +"testdata" + string(os.PathSeparator)
}