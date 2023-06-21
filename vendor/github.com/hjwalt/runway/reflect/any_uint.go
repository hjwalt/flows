package reflect

import (
	"math"
	"reflect"
	"strconv"

	"github.com/hjwalt/runway/logger"
	"go.uber.org/zap"
)

func GetUintBase(raw any, bitSize int) uint64 {
	input, isValid := GetValue(raw)
	if !isValid {
		return uint64(0)
	}
	var uintValue uint64
	switch input := input.(type) {
	case uint:
		uintValue = uint64(input)
	case uint8:
		uintValue = uint64(input)
	case uint16:
		uintValue = uint64(input)
	case uint32:
		uintValue = uint64(input)
	case uint64:
		uintValue = input
	case int, int8, int16, int32, int64:
		inputValue := GetIntBase(input, 64)
		uintValue = uint64(inputValue)
	case float32, float64:
		uintValue = uint64(math.Round(GetFloatBase(input, 64)))
	case bool:
		if input {
			uintValue = 1
		} else {
			uintValue = 0
		}
	case string:
		if input == "" {
			input = "0"
		}
		var err error
		uintValue, err = strconv.ParseUint(input, 10, bitSize)
		if err != nil {
			logger.WarnErr("string parse uint failed", err)
			return 0
		}
	case []byte:
		if len(input) == 0 {
			return 0
		}
		return Endian().Uint64(input)
	default:
		logger.Warn("conversion for uint type failed", zap.Any("type", reflect.TypeOf(input)), zap.Any("value", input))
		return 0
	}

	return uintValue
}

func GetUint(input any) uint {
	return uint(GetUintBase(input, 0))
}

func GetUint8(input any) uint8 {
	return uint8(GetUintBase(input, 8))
}

func GetUint16(input any) uint16 {
	return uint16(GetUintBase(input, 16))
}

func GetUint32(input any) uint32 {
	return uint32(GetUintBase(input, 32))
}

func GetUint64(input any) uint64 {
	return uint64(GetUintBase(input, 64))
}
