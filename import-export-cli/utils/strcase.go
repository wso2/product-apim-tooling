package utils

import "strings"

func ToPascalCase(s string) string {
	builder := strings.Builder{}
	tokens := strings.Split(strings.TrimSpace(s), " ")
	for _, token := range tokens {
		builder.WriteString(strings.Title(token))
	}
	return builder.String()
}
