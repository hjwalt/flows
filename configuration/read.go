package configuration

import (
	"os"

	"github.com/hjwalt/flows/format"
)

func Read[T any](file string, f format.Format[T]) (T, error) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		return f.Default(), err
	}
	return f.Unmarshal(bytes)
}
