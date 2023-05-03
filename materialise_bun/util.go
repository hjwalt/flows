package materialise_bun

import (
	"reflect"

	"github.com/gobeam/stringy"
)

var IgnoreColumns = map[string]bool{
	"base_model": true,
}

func Columns(val interface{}) ([]string, []string) {
	pks := make([]string, 0)
	cols := make([]string, 0)

	st := reflect.TypeOf(val)
	if st.Kind() == reflect.Pointer {
		st = st.Elem()
	}

	for i := 0; i < st.NumField(); i++ {
		if !st.Field(i).IsExported() {
			continue
		}
		fieldName := LowerSnakeCase(st.Field(i).Name)
		ispk := false
		if len(st.Field(i).Tag.Get("bun")) > 0 {
			tag := ParseTag(st.Field(i).Tag.Get("bun"))
			if len(tag.Name) > 0 {
				fieldName = tag.Name
			}
			_, ispk = tag.Options["pk"]
		}
		if _, shouldSkip := IgnoreColumns[fieldName]; shouldSkip {
			continue
		}
		if ispk {
			pks = append(pks, fieldName)
		} else {
			cols = append(cols, fieldName)
		}
	}

	return pks, cols
}

func ColumnsGeneric[T any](val T) ([]string, []string) {
	pks := make([]string, 0)
	cols := make([]string, 0)

	st := reflect.TypeOf(val)
	if st.Kind() == reflect.Pointer {
		st = st.Elem()
	}

	for i := 0; i < st.NumField(); i++ {
		if !st.Field(i).IsExported() {
			continue
		}
		fieldName := LowerSnakeCase(st.Field(i).Name)
		ispk := false
		if len(st.Field(i).Tag.Get("bun")) > 0 {
			tag := ParseTag(st.Field(i).Tag.Get("bun"))
			if len(tag.Name) > 0 {
				fieldName = tag.Name
			}
			_, ispk = tag.Options["pk"]
		}
		if _, shouldSkip := IgnoreColumns[fieldName]; shouldSkip {
			continue
		}
		if ispk {
			pks = append(pks, fieldName)
		} else {
			cols = append(cols, fieldName)
		}
	}

	return pks, cols
}

func ConflictSet(cols []string) string {
	conflictSet := ""
	for _, column := range cols {
		if len(conflictSet) != 0 {
			conflictSet += ","
		}
		conflictSet += column + " = EXCLUDED." + column
	}
	return conflictSet
}

func LowerSnakeCase(fieldName string) string {
	str := stringy.New(fieldName)
	snakeStr := str.SnakeCase("?", "")
	return snakeStr.ToLower()
}
