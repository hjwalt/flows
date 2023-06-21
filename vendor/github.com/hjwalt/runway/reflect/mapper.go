package reflect

import (
	"reflect"
	"time"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// constructor
func NewMapper(configurations ...func(*Mapper) *Mapper) *Mapper {
	defaultHandler := &Mapper{
		FieldSearch: DefaultStringSearch,
		Handler: map[reflect.Kind]func(interface{}, reflect.Type) reflect.Value{
			reflect.Int: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetInt(input))
			},
			reflect.Int8: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetInt8(input))
			},
			reflect.Int16: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetInt16(input))
			},
			reflect.Int32: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetInt32(input))
			},
			reflect.Int64: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetInt64(input))
			},
			reflect.Uint: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetUint(input))
			},
			reflect.Uint8: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetUint8(input))
			},
			reflect.Uint16: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetUint16(input))
			},
			reflect.Uint32: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetUint32(input))
			},
			reflect.Uint64: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetUint64(input))
			},
			reflect.Float32: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetFloat32(input))
			},
			reflect.Float64: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetFloat64(input))
			},
			reflect.Bool: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetBool(input))
			},
			reflect.String: func(input interface{}, sliceType reflect.Type) reflect.Value {
				return reflect.ValueOf(GetString(input))
			},
		},
	}

	defaultHandler.Handler[reflect.Slice] = defaultHandler.SliceHandler
	defaultHandler.Handler[reflect.Pointer] = defaultHandler.PointerHandler
	defaultHandler.Handler[reflect.Struct] = defaultHandler.StructHandler

	for _, configuration := range configurations {
		defaultHandler = configuration(defaultHandler)
	}

	return defaultHandler
}

func WithMapperFieldSearch(fieldSearch func(string) []string) func(*Mapper) *Mapper {
	return func(m *Mapper) *Mapper {
		m.FieldSearch = fieldSearch
		return m
	}
}

// implementation
type Mapper struct {
	Handler     map[reflect.Kind]func(interface{}, reflect.Type) reflect.Value
	FieldSearch func(string) []string
}

func (mapper *Mapper) Set(target any, source any) any {
	if IsNil(target) {
		return nil
	}

	targetValue := reflect.ValueOf(target)
	targetType := targetValue.Type()

	if valueHandler, handlerExist := mapper.Handler[targetType.Kind()]; handlerExist {
		value := valueHandler(source, targetType)
		if targetType.Kind() == reflect.Pointer {
			// Reflecting into pointer
			reflect.Indirect(targetValue).Set(reflect.Indirect(value))
			return target
		} else {
			return value.Interface()
		}
	} else {
		logger.Error("type missing from handler", zap.Any("type", targetType))
		return target
	}
}

func (mapper *Mapper) GetFieldValue(inputMap map[string]interface{}, fieldName string) (interface{}, bool) {
	fieldNamesToSearch := mapper.FieldSearch(fieldName)
	for _, nameEntry := range fieldNamesToSearch {
		if value, exist := inputMap[nameEntry]; exist {
			return value, exist
		}
	}
	return nil, false
}

func (mapper *Mapper) StructHandler(input interface{}, structType reflect.Type) reflect.Value {
	// TODO: standardised special case
	// Special case for timestamp to parse from string
	if structType.Kind() == reflect.Struct && structType == reflect.TypeOf(time.Time{}) {
		timestampValue := GetTimestampMs(input)
		if timestampValue != nil {
			return reflect.ValueOf(*timestampValue)
		}
	}

	// Default handler cases
	structValue := reflect.Indirect(reflect.New(structType)) // New creates pointer
	if input == nil {
		return structValue
	}

	// Reflect from map string interface
	if inputMap, reflectOk := input.(map[string]interface{}); reflectOk {
		for i := 0; i < structType.NumField(); i++ {
			if !structType.Field(i).IsExported() {
				continue
			}

			field := structType.Field(i)
			fieldValue := structValue.FieldByName(field.Name)
			if value, exist := mapper.GetFieldValue(inputMap, field.Name); exist {
				if valueHandler, handlerExist := mapper.Handler[field.Type.Kind()]; handlerExist {
					fieldValue.Set(valueHandler(value, field.Type))
				} else {
					logger.Infof("Unknown type %v", field.Type.Kind())
				}
			}
		}
	}
	// TODO: reflect from struct
	return structValue
}

func (mapper *Mapper) PointerHandler(input interface{}, pointerType reflect.Type) reflect.Value {
	// TODO: standardised special case
	// Special case for timestamp to parse from string
	if pointerType.Elem().Kind() == reflect.Struct && pointerType.Elem() == reflect.TypeOf(timestamppb.Timestamp{}) {
		timestampValue := GetProtoTimestampMs(input)
		if timestampValue != nil {
			return reflect.ValueOf(timestampValue)
		}
	}
	if pointerType.Elem().Kind() == reflect.Struct && pointerType.Elem() == reflect.TypeOf(time.Time{}) {
		timestampValue := GetTimestampMs(input)
		if timestampValue != nil {
			return reflect.ValueOf(timestampValue)
		}
	}

	// Default handler cases
	if input == nil {
		return reflect.Zero(pointerType)
	}
	pointedValue := reflect.New(pointerType.Elem())
	if valueHandler, handlerExist := mapper.Handler[pointerType.Elem().Kind()]; handlerExist {
		reflect.Indirect(pointedValue).Set(valueHandler(input, pointerType.Elem()))
	} else {
		logger.Infof("Unknown type %v for pointer", pointerType.Elem().Kind())
	}
	return pointedValue
}

func (mapper *Mapper) SliceHandler(input interface{}, sliceType reflect.Type) reflect.Value {
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)
	if input == nil {
		return sliceValue
	}

	if reflect.TypeOf(input).Kind() == reflect.Slice {
		inputReflectValue := reflect.ValueOf(input)
		for i := 0; i < inputReflectValue.Len(); i++ {
			v := inputReflectValue.Index(i).Interface()
			if valueHandler, handlerExist := mapper.Handler[sliceType.Elem().Kind()]; handlerExist {
				sliceValue = reflect.Append(sliceValue, valueHandler(v, sliceType.Elem()))
			} else {
				logger.Infof("Unknown type %v for slice handler", sliceType.Elem().Kind())
			}
		}
	}

	return sliceValue
}
