package reflect

import (
	"errors"
	"reflect"
)

func GetField(source any, field string) (any, error) {
	interfaceValue := reflect.ValueOf(source)

	nonPointerValue := interfaceValue
	if interfaceValue.Kind() == reflect.Pointer {
		if interfaceValue.IsNil() {
			return nil, nil
		}
		nonPointerValue = interfaceValue.Elem()
	}

	if nonPointerValue.Kind() != reflect.Struct {
		return nil, errors.New("source is not a struct or pointer to a struct")
	}

	fieldValue := nonPointerValue.FieldByName(field)
	if fieldValue.IsValid() {
		return fieldValue.Interface(), nil
	} else {
		return nil, errors.New("invalid field for struct")
	}
}
