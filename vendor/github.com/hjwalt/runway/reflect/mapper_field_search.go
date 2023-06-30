package reflect

import (
	"strings"

	"github.com/hjwalt/runway/stringy"
)

func DefaultStringSearch(fieldName string) []string {
	result := make([]string, 5)
	result[0] = fieldName
	result[1] = strings.ToLower(fieldName)
	result[2] = strings.ToUpper(fieldName)
	result[3] = stringy.ToLowerSnakeCase(fieldName)
	result[4] = stringy.ToUpperSnakeCase(fieldName)
	return result
}
