package reflect

import "reflect"

// returns indirected value (or the same value if it is not a pointer), value is not default value (nil)
func GetValue(input any) (any, bool) {
	if input == nil {
		return input, false
	}

	if !IsPointer(input) {
		return input, true
	}

	reflectedValue := reflect.Indirect(reflect.ValueOf(input))
	if !reflectedValue.IsValid() {
		return input, false
	}

	return reflectedValue.Interface(), true
}

func IsPointer(input any) bool {
	if input == nil {
		return true
	}
	if reflect.TypeOf(input).Kind() == reflect.Pointer {
		return true
	}
	return false
}

func IsNil(input any) bool {
	if input == nil {
		return true
	}
	if !IsPointer(input) {
		return false
	}
	return reflect.ValueOf(input).IsNil()
}
