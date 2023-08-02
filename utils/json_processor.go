package utils

import "strings"

func EscapeJsonChars(input string) string {
	input = strings.ReplaceAll(input, "\\", "\\\\")
	input = strings.ReplaceAll(input, "\"", "\\\"")
	input = strings.ReplaceAll(input, "\n", "\\n")
	input = strings.ReplaceAll(input, "\r", "\\r")
	input = strings.ReplaceAll(input, "\t", "\\t")
	input = strings.ReplaceAll(input, "\b", "\\b")
	input = strings.ReplaceAll(input, "\f", "\\f")
	return input
}
