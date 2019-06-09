package utils

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func Test_toPascalCase(t *testing.T) {
	assert.Equal(t, "MyPascalCaseText", ToPascalCase("     my pascal          case text  "))
}
