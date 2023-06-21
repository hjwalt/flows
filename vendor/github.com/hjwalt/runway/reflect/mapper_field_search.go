package reflect

import (
	"strings"

	"github.com/gobeam/stringy"
)

func DefaultStringSearch(fieldName string) []string {
	result := make([]string, 5)
	str := stringy.New(fieldName)
	snakeStr := str.SnakeCase("?", "")
	result[0] = fieldName
	result[1] = strings.ToLower(fieldName)
	result[2] = strings.ToUpper(fieldName)
	result[3] = snakeStr.ToLower()
	result[4] = snakeStr.ToUpper()
	return result
}
