package reflect

import "reflect"

func Construct[T any]() T {
	tType := reflect.TypeOf(new(T)).Elem()
	if tType.Kind() == reflect.Pointer {
		return reflect.New(tType.Elem()).Interface().(T)
	} else if tType.Kind() == reflect.Interface {
		return *new(T)
	} else if tType.Kind() == reflect.Func {
		return *new(T)
	} else {
		return reflect.Indirect(reflect.New(tType)).Interface().(T)
	}
}
