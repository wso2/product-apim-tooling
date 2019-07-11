package utils

import (
	"fmt"
	"os"
	"regexp"

	"github.com/hashicorp/go-multierror"
)

// Match for $VAR and capture VAR inside a group
var re = regexp.MustCompile(`\$(\w+)`)

// ErrRequiredEnvKeyMissing represents error used for indicate environment key missing
type ErrRequiredEnvKeyMissing struct {
	// Key is the missing entity
	Key string
}

func (e ErrRequiredEnvKeyMissing) Error() string {
	return fmt.Sprintf("%s is required, please set the environment variable", e.Key)
}

// EnvSubstitute substitutes variables from environment to the content. It uses regex to match variables and look up them in the
// environment before processing.
// returns an error if anything happen
func EnvSubstitute(content string) (string, error) {
	var errorResults error
	missingEnvKeys := false
	matches := re.FindAllStringSubmatch(content, -1) // matches is [][]string

	for _, match := range matches {
		Logln("Looking for:", match[0])
		if os.Getenv(match[1]) == "" {
			missingEnvKeys = true
			errorResults = multierror.Append(errorResults, &ErrRequiredEnvKeyMissing{Key: match[0]})
		}
	}

	if missingEnvKeys {
		return "", errorResults
	}

	expanded := os.ExpandEnv(content)
	return expanded, nil
}

