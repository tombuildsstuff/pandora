package utils

import (
	"strings"
)

func NormalizePropertyName(input string) string {
	return strings.Title(input)
}
