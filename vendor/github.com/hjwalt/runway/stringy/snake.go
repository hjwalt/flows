package stringy

import (
	"regexp"
	"strings"
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func ToLowerSnakeCase(fieldName string) string {
	snake := matchFirstCap.ReplaceAllString(fieldName, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToUpperSnakeCase(fieldName string) string {
	snake := matchFirstCap.ReplaceAllString(fieldName, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToUpper(snake)
}
