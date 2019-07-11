package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectEnvShouldFailWhenEnvNotPresent(t *testing.T) {
	data := `$MYVAR`
	str, err := EnvSubstitute(data)
	assert.Equal(t, "", str, "Should return empty string")
	assert.Error(t, err, "Should return an error")
}

func TestInjectEnvShouldPassWhenEnvPresents(t *testing.T) {
	data := `$MYVAR`
	_ = os.Setenv("MYVAR", "myval")
	str, err := EnvSubstitute(data)
	assert.Nil(t, err, "Error should be null")
	assert.Equal(t, "myval", str, "Should correctly replace environment variable")
}
